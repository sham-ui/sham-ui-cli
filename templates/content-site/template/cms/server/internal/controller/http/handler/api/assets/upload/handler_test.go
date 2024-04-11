package upload

import (
	"bytes"
	"cms/internal/controller/http/handler/assets"
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/logger"
	"cms/pkg/logger/testlogger"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func newUploadRequest(t *testing.T, fieldname, filename, content string) *http.Request {
	t.Helper()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filename)
	asserts.NoError(t, err)
	_, err = io.Copy(part, bytes.NewBufferString(content))
	asserts.NoError(t, err)
	err = writer.Close()
	asserts.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/upload-image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name                   string
		request                *http.Request
		assetsService          func(t mockConstructorTestingTNewMockAssetsService) *MockAssetsService
		expectedCode           int
		expectedBody           string
		expectedLoggerMessages []testlogger.Message
	}{
		{
			name:    "success",
			request: newUploadRequest(t, "file-0", "test.txt", "test"),
			assetsService: func(t mockConstructorTestingTNewMockAssetsService) *MockAssetsService {
				m := NewMockAssetsService(t)
				m.
					On("Save", mock.Anything, "test.txt", []byte("test")).
					Return(&model.Asset{
						Path:    "123.txt",
						Content: []byte("test"),
					}, nil).
					Once()
				return m
			},
			expectedCode:           http.StatusOK,
			expectedBody:           `{"result":[{"url":"/assets/123.txt","name":"test.txt","size":"4"}]}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "missing file",
			request:                httptest.NewRequest(http.MethodPost, "/upload-image", nil),
			assetsService:          NewMockAssetsService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Messages":["Fail parse multipart form"], "Status": "Bad Request"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:                   "invalid file key",
			request:                newUploadRequest(t, "file-1", "test.txt", "test"),
			assetsService:          NewMockAssetsService,
			expectedCode:           http.StatusBadRequest,
			expectedBody:           `{"Messages":["Fail get file"], "Status": "Bad Request"}`,
			expectedLoggerMessages: []testlogger.Message{},
		},
		{
			name:    "save fail",
			request: newUploadRequest(t, "file-0", "test.txt", "test"),
			assetsService: func(t mockConstructorTestingTNewMockAssetsService) *MockAssetsService {
				m := NewMockAssetsService(t)
				m.
					On("Save", mock.Anything, "test.txt", []byte("test")).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Messages":["internal server error"], "Status": "Internal Server Error"}`,
			expectedLoggerMessages: []testlogger.Message{
				{
					Message: "Fail save file",
					KeyValues: map[string]any{
						"error": "test",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			assets.Setup(router.PathPrefix("/assets").Subrouter(), nil)
			Setup(router, test.assetsService(t))
			log := testlogger.NewLogger()
			ctx := logger.Save(test.request.Context(), logr.New(log))
			req := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Action
			router.ServeHTTP(resp, req)

			// Assert
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
			asserts.JSONEquals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedLoggerMessages, log.Messages, "logger")
		})
	}
}
