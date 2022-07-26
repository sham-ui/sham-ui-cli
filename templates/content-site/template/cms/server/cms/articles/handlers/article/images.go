package article

import (
	"cms/config"
	"cms/core/handler"
	"cms/core/sessions"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"strings"
)

type ImagesHandler struct {
	fileServer http.Handler
}

type imageRequest struct {
	filePath string
}

func (h ImagesHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	path := mux.Vars(ctx.Request)["file"]
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	path = strings.ToLower(path)
	cannonicalName := strings.Replace(path, "\\", "/", -1)
	return imageRequest{
		filePath: cannonicalName,
	}, nil
}

func (h ImagesHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	req := data.(imageRequest)
	validation := handler.NewValidation()
	if 0 == len(req.filePath) {
		validation.AddError("empty path")
	}
	if "" == filepath.Ext(req.filePath) {
		validation.AddError("empty extension")
	}
	return validation, nil
}

func (h ImagesHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	req := data.(imageRequest)
	ctx.Request.URL.Path = req.filePath
	ctx.Response.Header().Set("Vary", "Accept-Encoding")
	ctx.Response.Header().Set("Cache-Control", "public, max-age=7776000")
	h.fileServer.ServeHTTP(ctx.Response, ctx.Request)
	return nil, nil
}

func NewImagesHandler(sessionStore *sessions.Store) http.HandlerFunc {
	h := handler.Create(
		ImagesHandler{
			fileServer: http.FileServer(http.Dir(config.Upload.Path)),
		},
		handler.WithOnlyForAuthenticated(sessionStore),
		handler.WithoutSerializeResultToJSON(),
	)
	return gziphandler.GzipHandler(h).ServeHTTP
}
