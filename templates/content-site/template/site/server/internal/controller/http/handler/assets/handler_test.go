package assets

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"site/internal/model"
	"site/pkg/asserts"
	"site/pkg/logger"

	"github.com/go-logr/logr/testr"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		name         string
		service      func(t mockConstructorTestingTNewMockService) *MockService
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			service: func(t mockConstructorTestingTNewMockService) *MockService {
				m := NewMockService(t)
				m.On("Asset", mock.Anything, "foo.txt").Return(&model.Asset{
					Path:    "foo",
					Content: []byte("123"),
				}, nil)
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/assets/foo.txt", nil),
			expectedCode: http.StatusOK,
			expectedBody: "123",
		},
		{
			name: "not found",
			service: func(t mockConstructorTestingTNewMockService) *MockService {
				m := NewMockService(t)
				m.On("Asset", mock.Anything, "foo.txt").Return(
					nil,
					model.AssetNotFoundError{Path: "foo.txt"},
				)
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/assets/foo.txt", nil),
			expectedCode: http.StatusNotFound,
			expectedBody: "{\"Status\":\"Not Found\",\"Messages\":[\"asset not found\"]}\n",
		},
		{
			name: "internal server error",
			service: func(t mockConstructorTestingTNewMockService) *MockService {
				m := NewMockService(t)
				m.On("Asset", mock.Anything, "foo.txt").Return(
					nil,
					errors.New("test"),
				)
				return m
			},
			request:      httptest.NewRequest(http.MethodGet, "/assets/foo.txt", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"Status\":\"Internal Server Error\",\"Messages\":[\"internal server error\"]}\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := mux.NewRouter()
			Setup(router, test.service(t))
			ctx := logger.Save(test.request.Context(), testr.New(t))
			request := test.request.WithContext(ctx)
			resp := httptest.NewRecorder()

			// Act
			router.ServeHTTP(resp, request)

			// Assert
			asserts.Equals(t, test.expectedBody, resp.Body.String(), "body")
			asserts.Equals(t, test.expectedCode, resp.Code, "code")
		})
	}
}
