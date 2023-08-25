package handlers

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"io"
	"io/fs"
	"net/http"
	"site/core/handler"
	"site/proto"
	"time"
)

type assetsHandler struct {
	cmsClient proto.CMSClient
}

type assetsRequest struct {
	filePath string
}

type fakeFileStats struct {
	name string
	size int64
}

func (f fakeFileStats) Name() string       { return f.name }
func (f fakeFileStats) Size() int64        { return f.size }
func (f fakeFileStats) Mode() fs.FileMode  { return fs.ModePerm }
func (f fakeFileStats) ModTime() time.Time { return time.Time{} }
func (f fakeFileStats) IsDir() bool        { return false }
func (f fakeFileStats) Sys() interface{}   { return nil }

var _ http.File = (*onceFile)(nil)

type onceFile struct {
	stats    fakeFileStats
	contents []byte
	offset   int64
}

func (o *onceFile) Close() error { return nil }

func (o *onceFile) Read(p []byte) (int, error) {
	if o.offset >= int64(len(o.contents)) {
		return 0, io.EOF
	}
	n := copy(p, o.contents[o.offset:])
	o.offset += int64(n)
	return n, nil
}

func (o *onceFile) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(o.offset) + offset
	case io.SeekEnd:
		abs = int64(len(o.contents)) + offset
	default:
		return 0, errors.New("onceFile.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("onceFile.Seek: negative position")
	}
	o.offset = abs
	return abs, nil
}

func (o *onceFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, fmt.Errorf("can't read dir")
}

func (o onceFile) Stat() (fs.FileInfo, error) {
	return o.stats, nil
}

var _ http.FileSystem = (*onceFileFS)(nil)

type onceFileFS struct {
	f *onceFile
}

func (o onceFileFS) Open(_ string) (http.File, error) {
	return o.f, nil
}

func (h *assetsHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	filePath, ok := mux.Vars(ctx.Request)["file"]
	if !ok {
		return nil, fmt.Errorf("can't extract file param")
	}
	return &assetsRequest{
		filePath: filePath,
	}, nil
}

func (h *assetsHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	params := data.(*assetsRequest)
	if "" == params.filePath {
		validation.AddError("file path must be not empty")
	}
	return validation, nil
}

func (h *assetsHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	params := data.(*assetsRequest)
	resp, err := h.cmsClient.Asset(ctx.Request.Context(), &proto.AssetRequest{
		Path: params.filePath,
	})
	if nil != err {
		return nil, fmt.Errorf("can't get file from api")
	}
	if _, ok := resp.Response.(*proto.AssetResponse_NotFound); ok {
		ctx.RespondWithError(http.StatusNotFound)
		return nil, nil
	}
	content := resp.GetFile()
	fileServer := http.FileServer(onceFileFS{
		f: &onceFile{
			stats: fakeFileStats{
				name: params.filePath,
				size: int64(len(content)),
			},
			contents: content,
		},
	})
	w := ctx.Response
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fileServer.ServeHTTP(w, ctx.Request)
	return nil, nil
}

func NewAssetsHandler(logger logr.Logger, cmsClient proto.CMSClient) http.HandlerFunc {
	h := &assetsHandler{
		cmsClient: cmsClient,
	}
	return handler.Create(logger, h, handler.WithoutSerializeResultToJSON())
}
