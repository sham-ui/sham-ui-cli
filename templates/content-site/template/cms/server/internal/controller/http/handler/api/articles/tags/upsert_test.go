package tags

import (
	"bytes"
	"cms/internal/model"
	"cms/pkg/asserts"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestExtractAndValidateData(t *testing.T) {
	testCases := []struct {
		name              string
		req               *http.Request
		slugger           func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService
		expectedTag       *model.Tag
		expectedDataValid bool
	}{
		{
			name: "success",
			req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{
				"name": "new tag"
			}`))),
			slugger: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.EXPECT().
					SlugifyTag(mock.Anything, "new tag").
					Return("new-tag").
					Once()
				return m
			},
			expectedTag: &model.Tag{ //nolint:exhaustruct
				Name: "new tag",
				Slug: "new-tag",
			},
			expectedDataValid: true,
		},
		{
			name:              "fail parse json",
			req:               httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(``))),
			slugger:           NewMockSlugifyService,
			expectedDataValid: false,
		},
		{
			name:              "empty name",
			req:               httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{}`))),
			slugger:           NewMockSlugifyService,
			expectedDataValid: false,
		},
		{
			name: "name trimmed",
			req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{
				"name": "   new tag    "
			}`))),
			slugger: func(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
				m := NewMockSlugifyService(t)
				m.EXPECT().
					SlugifyTag(mock.Anything, "new tag").
					Return("new-tag").
					Once()
				return m
			},
			expectedTag: &model.Tag{ //nolint:exhaustruct
				Name: "new tag",
				Slug: "new-tag",
			},
			expectedDataValid: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			rw := httptest.NewRecorder()

			// Action
			tag, valid := ExtractAndValidateData(test.slugger(t), rw, test.req)

			// Assert
			asserts.Equals(t, test.expectedTag, tag)
			asserts.Equals(t, test.expectedDataValid, valid)
		})
	}
}

func TestHandleError(t *testing.T) {
	testCases := []struct {
		name               string
		err                error
		expectedHandled    bool
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "slug already exists",
			err:                model.ErrTagSlugAlreadyExists,
			expectedHandled:    true,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"Status": "Bad Request", "Messages": ["Slug is already in use."]}`,
		},
		{
			name:               "name already exists",
			err:                model.ErrTagNameAlreadyExists,
			expectedHandled:    true,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"Status": "Bad Request", "Messages": ["Name is already in use."]}`,
		},
		{
			name:               "other error",
			err:                errors.New("test"),
			expectedHandled:    false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "err == nil",
			err:                nil,
			expectedHandled:    false,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()

			// Action
			handled := HandleError(test.err, rw, r)

			// Assert
			asserts.Equals(t, test.expectedHandled, handled)
			asserts.Equals(t, test.expectedStatusCode, rw.Code)
			if handled {
				asserts.JSONEquals(t, test.expectedBody, rw.Body.String())
			}
		})
	}
}
