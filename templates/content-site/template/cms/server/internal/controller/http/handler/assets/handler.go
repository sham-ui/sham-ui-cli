package assets

import (
	"cms/internal/controller/http/response"
	"cms/pkg/validation"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const RouteName = "assets"

const (
	FileKey  = "file"
	cacheTTL = 30 * 24 * time.Hour // 30 days
)

type handler struct {
	fsHandler http.Handler
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path, ok := mux.Vars(r)[FileKey]
	if !ok {
		response.BadRequest(rw, r, "Empty file")
		return
	}
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	path = strings.ToLower(path)
	path = strings.ReplaceAll(path, "\\", "/")

	valid := validation.New()
	if len(path) == 0 {
		valid.AddErrors("Empty file")
	}
	if filepath.Ext(path) == "" {
		valid.AddErrors("Empty extension")
	}
	if valid.HasErrors() {
		response.BadRequest(rw, r, valid.Errors...)
		return
	}

	r.URL.Path = path
	rw.Header().Set("Vary", "Accept-Encoding")
	rw.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(cacheTTL.Seconds())))

	h.fsHandler.ServeHTTP(rw, r)
}

func newHandler(fsHandler http.Handler) *handler {
	return &handler{
		fsHandler: fsHandler,
	}
}

func Setup(router *mux.Router, fs fs.FS) {
	fsHandler := http.FileServer(http.FS(fs))
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/{" + FileKey + "}").
		Handler(newHandler(fsHandler))
}
