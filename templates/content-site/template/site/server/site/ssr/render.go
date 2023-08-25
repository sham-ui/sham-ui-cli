package ssr

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/zealic/go2node"
	"os"
	"os/exec"
	"sync"
	"time"
)

//go:embed ssr.js
var script []byte

type Render interface {
	Start()
	Stop()
	render(ctx context.Context, req *nodejsRequest) ([]byte, error)
}

var _ Render = (*serverSideRender)(nil)

type serverSideRender struct {
	requests *requests
	cmd      *exec.Cmd
	channel  go2node.NodeChannel
	logger   logr.Logger

	mutex   sync.Mutex
	isRun   bool
	started chan bool
	done    chan bool
	restart chan bool
}

func (ssr *serverSideRender) Start() {
	cmd := exec.Command("node", "--version")
	resp, err := cmd.CombinedOutput()
	if nil != err {
		ssr.logger.Error(err, "Can't check node version")
		os.Exit(1)
	}
	ssr.logger.Info("Check node version", "version", string(resp))

	err = os.WriteFile("ssr.js", script, os.ModePerm)
	if nil != err {
		ssr.logger.Error(err, "Can't create ssr.js file")
		os.Exit(1)
	}
	go ssr.run()
	go ssr.autoRestart()
	<-ssr.started
}

func (ssr *serverSideRender) Stop() {
	ssr.mutex.Lock()
	if ssr.isRun {
		err := ssr.cmd.Process.Kill()
		if nil != err {
			ssr.logger.Error(err, "fail kill command process")
		}
		ssr.done <- true
		ssr.requests.ch <- false
	}
	ssr.mutex.Unlock()
}

func (ssr *serverSideRender) autoRestart() {
	for {
		select {
		case <-ssr.done:
			return
		default:
			ssr.mutex.Lock()
			cmd := ssr.cmd
			ssr.mutex.Unlock()
			if nil == cmd {
				ssr.initCommand()
				continue
			}
			_, err := cmd.Process.Wait()
			if nil != err {
				ssr.logger.Error(err, "SSR process stopped")
			}
			ssr.cmd = nil
			ssr.restart <- true
			time.Sleep(500 * time.Millisecond)
			ssr.logger.Info("Restart SSR")
		}
	}
}

func (ssr *serverSideRender) initCommand() {
	cmd := exec.Command("node", "ssr.js")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	channel, err := go2node.ExecNode(cmd)
	if nil != err {
		ssr.logger.Error(err, "can't create command")
		return
	}
	ssr.mutex.Lock()
	ssr.cmd = cmd
	ssr.channel = channel
	ssr.isRun = true
	ssr.mutex.Unlock()
	ssr.logger.Info("Start SSR sub process")
	ssr.started <- true
}

func (ssr *serverSideRender) run() {
	for ssr.requests.Next() {
		ssr.mutex.Lock()
		channel := ssr.channel
		cmd := ssr.cmd
		ssr.mutex.Unlock()

		resp := make(chan *go2node.NodeMessage, 1)

		go func(channel go2node.NodeChannel, resp chan *go2node.NodeMessage) {
			msg, err := channel.Read()
			defer close(resp)
			if nil != err {
				ssr.logger.Error(err, "fail read from ssr chanel")
				return
			}
			resp <- msg
		}(channel, resp)

		select {
		case <-ssr.restart:
			cmd.Process.Kill()
			<-ssr.started
			continue
		case msg := <-resp:
			var resp nodejsResponse
			if err := json.Unmarshal(msg.Message, &resp); nil != err {
				ssr.logger.Error(err, "unmarshal ssr response")
				continue
			}
			ssr.requests.publish(&resp)
		}

	}
}

func (ssr *serverSideRender) render(ctx context.Context, req *nodejsRequest) ([]byte, error) {
	data, err := json.Marshal(req)
	if nil != err {
		return nil, fmt.Errorf("can't marshal ssr request: %w", err)
	}
	ch := ssr.requests.subscribe(req)
	defer ssr.requests.unsubscribe(req.ID)
	err = ssr.channel.Write(&go2node.NodeMessage{
		Message: data,
	})
	if nil != err {
		return nil, fmt.Errorf("can't write message to channel: %w", err)
	}
	select {
	case resp := <-ch:
		if "" != resp.Error {
			return nil, fmt.Errorf("ssr respond with error: %s", resp.Error)
		}
		return []byte(resp.HTML), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewServerSideRender(logger logr.Logger) Render {
	return &serverSideRender{
		logger:   logger.WithName("ssr"),
		requests: newRequests(100),
		started:  make(chan bool, 1),
		done:     make(chan bool, 1),
		restart:  make(chan bool, 1),
	}
}
