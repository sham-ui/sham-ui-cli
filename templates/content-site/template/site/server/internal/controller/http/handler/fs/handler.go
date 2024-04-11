package fs

import (
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"site/internal/controller/http/response"
	"site/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "fs"

type handler struct {
	fileSystem  fs.FS
	fileServer  http.Handler
	htmlHandler http.Handler
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := logger.Load(r.Context())

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		log.Error(err, "can't get absolute path")
		response.BadRequest(rw, r, "can't get absolute path")
		return
	}

	// strip leading "/"
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// check whether a file exists at the given path
	cannonicalName := strings.ReplaceAll(path, "\\", "/")
	if cannonicalName == "" {
		h.htmlHandler.ServeHTTP(rw, r)
		return
	}

	file, err := h.fileSystem.Open(cannonicalName)
	if errors.Is(err, fs.ErrNotExist) {
		// file does not exist, serve as html
		h.htmlHandler.ServeHTTP(rw, r)
		return
	}
	if err != nil {
		log.Error(err, "can't open file", "path", cannonicalName)
		response.InternalServerError(rw, r)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if nil != err {
		log.Error(err, "can't get stat for file", "path", cannonicalName)
		response.InternalServerError(rw, r)
		return
	}

	if stat.IsDir() {
		// file is dir with assets, serve as html
		h.htmlHandler.ServeHTTP(rw, r)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	rw.Header().Set("Vary", "Accept-Encoding")
	rw.Header().Set("Cache-Control", "public, max-age=86400")
	h.fileServer.ServeHTTP(rw, r)
}

func newHandler(fileSystem fs.FS, htmlHandler http.Handler) *handler {
	return &handler{
		fileSystem:  fileSystem,
		fileServer:  http.FileServer(http.FS(fileSystem)),
		htmlHandler: htmlHandler,
	}
}

func Setup(router *mux.Router, fileSystem fs.FS, htmlHandler http.Handler) {
	router.
		Name(RouteName).
		PathPrefix("/").
		Handler(newHandler(fileSystem, htmlHandler))
}
