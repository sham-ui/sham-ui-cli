package upload

import (
	"cms/internal/controller/http/handler/assets"
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const RouteName = "api.assets.upload"

const (
	maxUploadSizeBytes = 32 << 20 // 32 MB
	fileKey            = "file-0"
)

type (
	uploadResponse struct {
		Result []uploadFileResponse `json:"result"`
	}

	uploadFileResponse struct {
		URL  string `json:"url"`
		Name string `json:"name"`
		Size string `json:"size"`
	}
)

type handler struct {
	router        *mux.Router
	assetsService AssetsService
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(maxUploadSizeBytes); err != nil {
		response.BadRequest(rw, r, "Fail parse multipart form")
		return
	}
	filePath, header, err := r.FormFile(fileKey)
	if err != nil {
		response.BadRequest(rw, r, "Fail get file")
		return
	}
	fileContent, err := io.ReadAll(filePath)
	if err != nil {
		response.BadRequest(rw, r, "Fail read file")
		return
	}

	asset, err := h.assetsService.Save(ctx, header.Filename, fileContent)
	if err != nil {
		logger.Load(ctx).Error(err, "Fail save file")
		response.InternalServerError(rw, r)
		return
	}

	url, err := h.router.Get(assets.RouteName).URL(assets.FileKey, asset.Path)
	if err != nil {
		logger.Load(ctx).Error(err, "Fail build asset route")
		response.InternalServerError(rw, r)
		return
	}

	response.JSON(rw, r, http.StatusOK, uploadResponse{
		Result: []uploadFileResponse{
			{
				URL:  url.String(),
				Name: header.Filename,
				Size: strconv.FormatInt(header.Size, 10),
			},
		},
	})
}

func newHandler(assetsService AssetsService, router *mux.Router) *handler {
	return &handler{
		router:        router,
		assetsService: assetsService,
	}
}

func Setup(router *mux.Router, assetsService AssetsService) {
	router.
		Name(RouteName).
		Methods(http.MethodPost).
		Path("/upload-image").
		Handler(newHandler(assetsService, router))
}
