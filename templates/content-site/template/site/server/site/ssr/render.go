package ssr

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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

	mutex   sync.Mutex
	isRun   bool
	started chan bool
	done    chan bool
}

func (ssr *serverSideRender) Start() {
	cmd := exec.Command("node", "--version")
	resp, err := cmd.CombinedOutput()
	if nil != err {
		log.WithError(err).Fatalf("Can't check node version")
	}
	log.WithField("version", string(resp)).Info("Check node version")

	err = os.WriteFile("ssr.js", script, os.ModePerm)
	if nil != err {
		log.WithError(err).Fatalf("Can't create ssr.js file")
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
			log.WithError(err).Errorf("fail kill command process")
		}
		ssr.done <- true
		ssr.requests.ch <- false
	}
	ssr.mutex.Unlock()
}

func (ssr *serverSideRender) autoRestart() {
	needPublish := true
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
			if needPublish {
				needPublish = false
				ssr.started <- true
			}
			_, err := cmd.Process.Wait()
			if nil != err {
				log.WithError(err).Error("SSR process stopped")
			}
			ssr.cmd = nil
			ssr.channel = nil
			time.Sleep(500 * time.Millisecond)
			log.Info("Restart SSR")
			ssr.initCommand()
		}
	}
}

func (ssr *serverSideRender) initCommand() {
	cmd := exec.Command("node", "ssr.js")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	channel, err := go2node.ExecNode(cmd)
	if nil != err {
		log.WithError(err).Errorf("can't create command")
		return
	}
	ssr.mutex.Lock()
	ssr.cmd = cmd
	ssr.channel = channel
	ssr.isRun = true
	ssr.mutex.Unlock()
	log.Info("Start SSR sub process")
}

func (ssr *serverSideRender) run() {
	for ssr.requests.Next() {
		ssr.mutex.Lock()
		channel := ssr.channel
		ssr.mutex.Unlock()
		msg, err := channel.Read()
		if nil != err {
			log.WithError(err).Errorf("read from ssr chanel")
			continue
		}
		var resp nodejsResponse
		if err := json.Unmarshal(msg.Message, &resp); nil != err {
			log.Errorf("unmarshal ssr response: %s", err)
			continue
		}
		ssr.requests.publish(&resp)
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

func NewServerSideRender() Render {
	return &serverSideRender{
		requests: newRequests(100),
		started:  make(chan bool, 1),
		done:     make(chan bool, 1),
	}
}
