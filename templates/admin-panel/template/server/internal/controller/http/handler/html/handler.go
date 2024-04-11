package html

import (
	"bytes"
	"net/http"
	"time"
)

type handler struct {
	htmlFile []byte
	modTime  time.Time
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.ServeContent(rw, r, "index.html", h.modTime, bytes.NewReader(h.htmlFile))
}

func New(html []byte) *handler {
	return &handler{
		htmlFile: html,
		modTime:  time.Now(),
	}
}
