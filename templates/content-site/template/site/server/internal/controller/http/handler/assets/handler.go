package assets

import (
	"errors"
	"net/http"

	"site/internal/controller/http/response"
	"site/internal/model"
	"site/pkg/logger"
	"site/pkg/one_file_fs"

	"github.com/gorilla/mux"
)

const paramKey = "file"

type handler struct {
	service Service
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
	http.FileServer(one_file_fs.New(asset.Path, asset.Content)).ServeHTTP(rw, r)
	rw.Header().Set("Cache-Control", "public, max-age=7776000")
}

func newHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func Setup(router *mux.Router, service Service) {
	router.
		Methods("GET").
		Path("/assets/{" + paramKey + "}").
		Handler(newHandler(service))
}
