package ssr

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
	"github.com/zealic/go2node"
	"net/http"
	"os"
	"os/exec"
	"site/config"
	"strconv"
)

//go:embed ssr.js
var script []byte

type ServerSideRender struct {
	cmd      *exec.Cmd
	channel  go2node.NodeChannel
	requests *requests
	apiURL   string
	done     chan bool
}

func (ssr *ServerSideRender) Stop() {
	ssr.done <- true
	ssr.cmd.Process.Kill()
}

func (ssr *ServerSideRender) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	id, err := gonanoid.New()
	if nil != err {
		log.WithError(err).Error("can't generate nanoid")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cookiesString string
	cookies := r.Cookies()
	if len(cookies) > 0 {
		cookiesString = cookies[0].String()
		if len(cookies) > 1 {
			for _, item := range cookies[1:] {
				cookiesString += "; " + item.String()
			}
		}
	}

	origin := ssr.getOrigin(r)
	resp, err := ssr.render(&nodejsRequest{
		ID:      id,
		URL:     origin + r.URL.Path + "?" + r.URL.RawQuery,
		Origin:  origin,
		API:     ssr.apiURL,
		Cookies: cookiesString,
	})
	if nil != err {
		log.WithError(err).Errorf("can't ssr")
		http.Error(w, "SSR error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if nil != err {
		log.WithError(err).Errorf("can't write ssr response")
	}
}

func (ssr *ServerSideRender) getOrigin(r *http.Request) string {
	var url string
	if nil == r.TLS {
		url += "http://"
	} else {
		url += "https://"
	}
	return url + r.Host
}

func (ssr *ServerSideRender) render(req *nodejsRequest) ([]byte, error) {
	data, err := json.Marshal(req)
	if nil != err {
		return nil, fmt.Errorf("can't marshal ssr request: %w", err)
	}
	ch := make(chan *nodejsResponse, 1)
	ssr.requests.subscribe(req, ch)
	err = ssr.channel.Write(&go2node.NodeMessage{
		Message: data,
	})
	if nil != err {
		return nil, fmt.Errorf("can't write message to channel: %w", err)
	}
	resp := <-ch
	if "" != resp.Error {
		return nil, fmt.Errorf("ssr respond with error: %s", resp.Error)
	}
	return []byte(resp.HTML), nil
}

func (ssr *ServerSideRender) run() {
	for {
		select {
		case <-ssr.done:
			return
		default:
			resp, err := ssr.channel.Read()
			if nil != err {
				log.WithError(err).Errorf("read from ssr chanel")
				return
			}
			var msg nodejsResponse
			if err := json.Unmarshal(resp.Message, &msg); nil != err {
				log.Errorf("unmarshal ssr response: %s", err)
				return
			}
			ssr.requests.publish(&msg)
		}
	}
}

func NewServerSideRender() *ServerSideRender {
	cmd := exec.Command("node", "--version")
	resp, err := cmd.CombinedOutput()
	if nil != err {
		log.WithError(err).Fatalf("can't check node version")
	}
	log.WithField("version", string(resp)).Info("check node version")

	err = os.WriteFile("ssr.js", script, os.ModePerm)
	if nil != err {
		log.WithError(err).Fatalf("can't create ssr.js file")
	}

	cmd = exec.Command("node", "ssr.js")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	channel, err := go2node.ExecNode(cmd)
	if nil != err {
		log.WithError(err).Fatalf("can't create go2node channel")
	}

	go func() {
		_, err := cmd.Process.Wait()
		if nil != err {
			log.WithError(err).Error("ssr process fail")
		}
	}()

	s := &ServerSideRender{
		cmd:      cmd,
		channel:  channel,
		requests: newRequests(),
		apiURL:   "http://localhost:" + strconv.Itoa(config.Server.Port) + "/api/",
		done:     make(chan bool, 1),
	}

	go s.run()
	log.Info("Start SSR sub process")

	return s
}
