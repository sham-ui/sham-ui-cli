package assets

import (
	"cms/core/sessions"
	"embed"
	"errors"
	"github.com/NYTimes/gziphandler"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed files/*
var Assets embed.FS

type Handler struct {
	sessionStore *sessions.Store
	fsys         fs.FS
	fileServer   http.Handler
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// strip leading "/"
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// check whether a file exists at the given path
	cannonicalName := strings.Replace(path, "\\", "/", -1)
	file, err := h.fsys.Open(cannonicalName)
	if nil != file {
		defer file.Close()
	}
	if errors.Is(err, fs.ErrNotExist) {

		// file does not exist, serve index.html
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		content, err := Assets.ReadFile("files/index.html")
		if nil != err {
			log.Errorf("can't get asset index.html: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(content)
		if nil != err {
			log.Errorf("can't write index.html asset: %s", err)
		}
	} else if strings.HasPrefix(cannonicalName, "su_") {
		session, err := h.sessionStore.GetSession(r)
		if nil != err {
			log.Errorf("can't get session: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if nil == session {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else if session.IsSuperuser {
			h.fileServe(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	} else {

		// otherwise, use http.FileServer to serve the static dir
		h.fileServe(w, r)
	}
}

func (h Handler) fileServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	h.fileServer.ServeHTTP(w, r)
}

func NewHandler(sessionStore *sessions.Store) http.Handler {
	fsys, err := fs.Sub(Assets, "files")
	if nil != err {
		log.WithError(err).Fatal("can' get file system")
	}
	return gziphandler.GzipHandler(Handler{
		sessionStore: sessionStore,
		fsys:         fsys,
		fileServer:   http.FileServer(http.FS(fsys)),
	})
}
