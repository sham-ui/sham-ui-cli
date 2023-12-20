package ssr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"

	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-logr/logr"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/zealic/go2node"
)

const (
	scopeName        = "service.ssr."
	requestPoolSize  = 100
	maxNodeJSTimeout = 30 * time.Second
)

type go2NodeResponse struct {
	success *go2node.NodeMessage
	error   error
}

type service struct {
	logger            logr.Logger
	tracer            trace.Tracer
	textMapPropagator propagation.TextMapPropagator
	script            []byte
	apiURL            string
	requests          *requests
	cmd               *exec.Cmd
	channel           go2node.NodeChannel
	isRun             bool
	mutex             sync.Mutex
	done              chan bool
	restart           chan bool
}

func (srv *service) String() string {
	return "service.ssr"
}

func (srv *service) GracefulShutdown(_ context.Context) error {
	srv.Stop()
	return nil
}

func (srv *service) Start() error {
	cmd := exec.Command("node", "--version")
	resp, err := cmd.CombinedOutput()
	if nil != err {
		return fmt.Errorf("check node version: %w", err)
	}
	srv.logger.Info("check node version", "version", string(resp))
	if err = os.WriteFile("ssr.js", srv.script, os.ModePerm); err != nil {
		return fmt.Errorf("create ssr.js file: %w", err)
	}
	if err := srv.initCommand(); err != nil {
		return fmt.Errorf("init command: %w", err)
	}
	go srv.run()
	go srv.autoRestart()
	return nil
}

func (srv *service) Stop() {
	srv.mutex.Lock()
	if srv.isRun {
		if err := srv.cmd.Process.Kill(); err != nil {
			srv.logger.Error(err, "fail kill command process")
		}
		srv.done <- true
		srv.requests.ch <- false
	}
	srv.mutex.Unlock()
}

func (srv *service) Render(ctx context.Context, url *url.URL, cookies []*http.Cookie) ([]byte, error) {
	const op = scopeName + "Render"

	ctx, span := srv.tracer.Start(ctx, op)
	defer span.End()

	id, err := gonanoid.New()
	if err != nil {
		return nil, fmt.Errorf("generate nanoid: %w", err)
	}

	var cookiesString string
	if len(cookies) > 0 {
		cookiesString = cookies[0].String()
		if len(cookies) > 1 {
			for _, item := range cookies[1:] {
				cookiesString += "; " + item.String()
			}
		}
	}

	headers := make(propagation.MapCarrier)
	srv.textMapPropagator.Inject(ctx, headers)

	origin := url.Scheme + "://" + url.Host
	req := &nodejsRequest{
		ID:      id,
		URL:     origin + url.Path + "?" + url.RawQuery,
		Origin:  origin,
		API:     srv.apiURL,
		Cookies: cookiesString,
		Headers: headers,
	}

	subscribe := srv.requests.subscribe(ctx, req)
	defer srv.requests.unsubscribe(req.ID)
	if err := srv.sendToNodeJs(ctx, req); err != nil {
		return nil, err
	}

	select {
	case resp := <-subscribe:
		if resp.Error != "" {
			span.SetStatus(codes.Error, "ssr respond with error")
			return nil, NewServerSideRenderError(resp.Error)
		}
		span.SetStatus(codes.Ok, "ssr respond with success")
		return []byte(resp.HTML), nil
	case <-ctx.Done():
		span.SetStatus(codes.Error, "ssr request timeout")
		return nil, ctx.Err()
	}
}

func (srv *service) sendToNodeJs(ctx context.Context, req *nodejsRequest) error {
	const op = scopeName + "sendToNodeJs"

	_, span := srv.tracer.Start(ctx, op)
	defer span.End()

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	if err := srv.channel.Write(&go2node.NodeMessage{Message: data}); err != nil {
		return NewServerSideRenderRequestError(err)
	}
	return nil
}

func (srv *service) initCommand() error {
	cmd := exec.Command("node", "ssr.js")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	channel, err := go2node.ExecNode(cmd)
	if err != nil {
		return fmt.Errorf("create go2node channel: %w", err)
	}
	srv.mutex.Lock()
	srv.cmd = cmd
	srv.channel = channel
	srv.isRun = true
	srv.mutex.Unlock()
	srv.logger.Info("start ssr sub process")
	return nil
}

func (srv *service) run() {
	for srv.requests.Next() {
		ctx, cancel := context.WithTimeout(context.Background(), maxNodeJSTimeout)
		if err := srv.processRequest(ctx); err != nil {
			srv.logger.Error(err, "fail process ssr request")
		}
		cancel()
	}
}

func (srv *service) processRequest(ctx context.Context) error {
	const op = scopeName + "processRequest"

	ctx, span := srv.tracer.Start(ctx, op)
	defer span.End()

	srv.mutex.Lock()
	channel := srv.channel
	cmd := srv.cmd
	srv.mutex.Unlock()
	if cmd == nil || channel == nil {
		span.SetStatus(codes.Error, "no ssr process")
		return errSSRProcessStopped
	}

	go2Node := make(chan go2NodeResponse, 1)

	go func() {
		msg, err := channel.Read()
		go2Node <- go2NodeResponse{
			success: msg,
			error:   err,
		}
		close(go2Node)
	}()

	msg, err := srv.waitResponse(ctx, cmd, go2Node)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var nr nodejsResponse
	if err := json.Unmarshal(msg.Message, &nr); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	srv.requests.publish(&nr)

	return nil
}

func (srv *service) waitResponse(
	ctx context.Context,
	cmd *exec.Cmd,
	go2Node <-chan go2NodeResponse,
) (*go2node.NodeMessage, error) {
	const op = scopeName + "waitResponse"

	ctx, span := srv.tracer.Start(ctx, op)
	defer span.End()

	select {
	case <-srv.restart:
		err := cmd.Process.Kill()
		if err == nil || errors.Is(err, os.ErrProcessDone) {
			return nil, errSSRProcessStopped
		}
		return nil, fmt.Errorf("kill ssr process: %w", err)
	case resp, ok := <-go2Node:
		if !ok {
			return nil, fmt.Errorf("channel closed: %w", ctx.Err())
		}
		if resp.error != nil {
			return nil, resp.error
		}
		return resp.success, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (srv *service) autoRestart() {
	for {
		select {
		case <-srv.done:
			return
		default:
			srv.mutex.Lock()
			cmd := srv.cmd
			srv.mutex.Unlock()

			if cmd == nil {
				if err := srv.initCommand(); err != nil {
					srv.logger.Error(err, "can't init command")
				}
				continue
			}

			if _, err := cmd.Process.Wait(); err != nil {
				srv.logger.Error(err, "ssr process stopped")
			}

			srv.mutex.Lock()
			srv.cmd = nil
			srv.isRun = false
			srv.mutex.Unlock()
			srv.restart <- true
			srv.logger.Info("restart ssr")
		}
	}
}

func New(
	logger logr.Logger,
	tracing trace.TracerProvider,
	textMapPropagator propagation.TextMapPropagator,
	apiURL string,
	script []byte,
) *service {
	return &service{
		tracer:            tracing.Tracer(scopeName),
		textMapPropagator: textMapPropagator,
		logger:            logger.WithName(scopeName),
		script:            script,
		apiURL:            apiURL,
		requests:          newRequests(requestPoolSize),
		done:              make(chan bool, 1),
		restart:           make(chan bool, 1),
		isRun:             false,
		cmd:               nil,
		channel:           nil,
		mutex:             sync.Mutex{},
	}
}
