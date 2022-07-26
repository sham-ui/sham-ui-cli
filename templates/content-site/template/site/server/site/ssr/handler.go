package ssr

import (
	"context"
	"github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"site/config"
	"strconv"
	"time"
)

type server struct {
	apiURL string
	render Render
}

func (ssr *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

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
	resp, err := ssr.render.render(ctx, &nodejsRequest{
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

func (ssr *server) getOrigin(r *http.Request) string {
	var url string
	if nil == r.TLS {
		url += "http://"
	} else {
		url += "https://"
	}
	return url + r.Host
}

func NewServer(render Render) http.Handler {
	return &server{
		apiURL: "http://localhost:" + strconv.Itoa(config.Server.Port) + "/api/",
		render: render,
	}
}
