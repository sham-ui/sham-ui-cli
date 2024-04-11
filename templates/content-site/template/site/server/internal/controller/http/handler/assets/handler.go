package assets

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"site/internal/controller/http/response"
	"site/internal/model"
	"site/pkg/logger"

	"github.com/gorilla/mux"
)

const RouteName = "assets"

const paramKey = "file"

type handler struct {
	service Service
	modTime time.Time
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log := logger.Load(r.Context())
	filePath := mux.Vars(r)[paramKey]
	log = log.WithValues("file", filePath)
	asset, err := h.service.Asset(r.Context(), filePath)
	if errors.Is(err, model.AssetNotFoundError{}) {
		response.NotFound(rw, r, "asset not found")
		return
	}
	if err != nil {
		response.InternalServerError(rw, r)
		log.Error(err, "can't get file")
		return
	}
	http.ServeContent(rw, r, asset.Path, h.modTime, bytes.NewReader(asset.Content))
	rw.Header().Set("Cache-Control", "public, max-age=7776000")
}

func newHandler(service Service, modTime time.Time) *handler {
	return &handler{
		service: service,
		modTime: modTime,
	}
}

func Setup(router *mux.Router, service Service) {
	router.
		Name(RouteName).
		Methods("GET").
		Path("/assets/{" + paramKey + "}").
		Handler(newHandler(service, time.Now()))
}
