package article

import (
	"cms/config"
	"cms/core/handler"
	"cms/core/sessions"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type uploadFileRequest struct {
	file   multipart.File
	header *multipart.FileHeader
}

type uploadResponse struct {
	ErrorMessage string               `json:"errorMessage"`
	Result       []uploadFileResponse `json:"result"`
}

type uploadFileResponse struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size string `json:"size"`
}

type uploadHandler struct {
}

func (h *uploadHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if nil != err {
		return nil, fmt.Errorf("parse multipart form: %s", err)
	}
	file, header, err := ctx.Request.FormFile("file-0")
	if nil != err {
		return nil, fmt.Errorf("get form file: %s", err)
	}
	return &uploadFileRequest{
		file:   file,
		header: header,
	}, nil
}

func (h *uploadHandler) Validate(ctx *handler.Context, data interface{}) (*handler.Validation, error) {
	return handler.NewValidation(), nil
}

func (h *uploadHandler) Process(_ *handler.Context, data interface{}) (interface{}, error) {
	req := data.(*uploadFileRequest)
	fi, err := io.ReadAll(req.file)
	if nil != err {
		return nil, fmt.Errorf("read file: %s", err)
	}
	fileName := h.getUniqFileName(fi, req.header.Filename)
	destPath := path.Join(config.Upload.Path, fileName)
	fileExists, err := h.checkFileExists(destPath)
	if nil != err {
		return nil, fmt.Errorf("check file exists: %s", err)
	}
	if fileExists {
		return h.buildFileResponse(fileName, req.header.Size), nil
	}
	fo, err := os.Create(destPath)
	if nil != err {
		return nil, fmt.Errorf("create file: %s", err)
	}
	defer fo.Close()
	n, err := fo.Write(fi)
	if nil != err {
		return nil, fmt.Errorf("write file: %s", err)
	}
	return h.buildFileResponse(fileName, int64(n)), nil
}

func (h *uploadHandler) getUniqFileName(content []byte, fileName string) string {
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:]) + strings.ToLower(filepath.Ext(fileName))
}

func (h *uploadHandler) checkFileExists(destPath string) (bool, error) {
	_, err := os.Stat(destPath)
	if nil == err {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (h *uploadHandler) buildFileResponse(fileName string, size int64) *uploadResponse {
	return &uploadResponse{
		Result: []uploadFileResponse{
			{
				URL:  "/assets/" + fileName,
				Name: fileName,
				Size: strconv.FormatInt(size, 10),
			},
		},
	}
}

func NewUploadHandler(sessionsStore *sessions.Store) http.HandlerFunc {
	h := &uploadHandler{}
	return handler.Create(h, handler.WithOnlyForAuthenticated(sessionsStore))
}
