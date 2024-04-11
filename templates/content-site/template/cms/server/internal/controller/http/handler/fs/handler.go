package fs

import (
	"cms/internal/controller/http/request"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"cms/internal/controller/http/response"
	"cms/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "fs"

const (
	cacheTTL             = 24 * time.Hour
	superuserFilesPrefix = "su_"
)

type handler struct {
	fileSystem  fs.FS
	fileServer  http.Handler
	htmlHandler http.Handler
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) { //nolint:funlen,cyclop
	ctx := r.Context()
	log := logger.Load(ctx)

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

	if strings.HasPrefix(cannonicalName, superuserFilesPrefix) {
		sess, ok := request.SessionFromContext(ctx)
		switch {
		case !ok || sess == nil:
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		case !sess.IsSuperuser:
			http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}

	file, err := h.fileSystem.Open(cannonicalName)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		// file does not exist, serve as html
		h.htmlHandler.ServeHTTP(rw, r)
		return
	case err != nil:
		log.Error(err, "can't open file", "path", cannonicalName)
		response.InternalServerError(rw, r)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
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
	rw.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(cacheTTL.Seconds())))
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
		Methods(http.MethodGet).
		Handler(newHandler(fileSystem, htmlHandler))
}
