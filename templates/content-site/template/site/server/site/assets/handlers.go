package assets

import (
	"embed"
	"errors"
	"github.com/NYTimes/gziphandler"
	"github.com/go-logr/logr"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"site/ssr"
	"strings"
)

//go:embed files/*
var Assets embed.FS

type Handler struct {
	ssr        http.Handler
	fsys       fs.FS
	fileServer http.Handler
	logger     logr.Logger
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

	fsys, err := fs.Sub(Assets, "files")
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check whether a file exists at the given path
	cannonicalName := strings.Replace(path, "\\", "/", -1)
	if "" == cannonicalName {
		h.ssr.ServeHTTP(w, r)
		return
	}

	file, err := fsys.Open(cannonicalName)
	if nil != file {
		defer file.Close()
	}
	if errors.Is(err, fs.ErrNotExist) {
		// file does not exist, serve SSR
		h.ssr.ServeHTTP(w, r)
	} else if err != nil {
		h.logger.Error(err, "can't open file", "path", cannonicalName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		stat, err := file.Stat()
		if nil != err {
			h.logger.Error(err, "can't get stat for file", "path", cannonicalName)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if stat.IsDir() {
			// file is dir with assets, serve SSR
			h.ssr.ServeHTTP(w, r)
			return
		}

		// otherwise, use http.FileServer to serve the static dir
		h.fileServe(w, r)
	}
}

func (h Handler) fileServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	h.fileServer.ServeHTTP(w, r)
}

func NewHandler(logger logr.Logger, render ssr.Render) http.Handler {
	assetsLogger := logger.WithName("assets")
	fsys, err := fs.Sub(Assets, "files")
	if nil != err {
		assetsLogger.Error(err, "can't get file system")
		os.Exit(1)
	}
	return gziphandler.GzipHandler(Handler{
		ssr:        ssr.NewServer(logger, render),
		fsys:       fsys,
		fileServer: http.FileServer(http.FS(fsys)),
		logger:     assetsLogger,
	})
}
