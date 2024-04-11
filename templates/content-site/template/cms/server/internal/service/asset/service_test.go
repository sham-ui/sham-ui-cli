package asset

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"context"
	"errors"
	"testing"
)

func TestService_Save_Existed(t *testing.T) {
	// Arrange
	srv := New("./testdata")

	// Action
	asset, err := srv.Save(context.Background(), "existed.txt", []byte("test"))

	// Arrange
	asserts.Equals(t, &model.Asset{
		Path:    "098f6bcd4621d373cade4e832627b4f6.txt",
		Content: []byte("test"),
	}, asset)
	asserts.NoError(t, err)
}

func TestService_Save_New(t *testing.T) {
	// Arrange
	uploadPath := t.TempDir()
	srv := New(uploadPath)

	// Action
	asset, err := srv.Save(context.Background(), "new.txt", []byte("test"))

	// Arrange
	asserts.Equals(t, &model.Asset{
		Path:    "098f6bcd4621d373cade4e832627b4f6.txt",
		Content: []byte("test"),
	}, asset)
	asserts.NoError(t, err)
}

func TestService_ReadFile(t *testing.T) {
	// Arrange
	srv := New("./testdata")

	testCases := []struct {
		path             string
		expectedResponse []byte
		expectedError    error
	}{
		{
			path:             "test.txt",
			expectedResponse: []byte("test"),
		},
		{
			path:             "not-found.txt",
			expectedResponse: nil,
			expectedError:    errors.New("asset not found: testdata/not-found.txt"),
		},
	}
	for _, test := range testCases {
		t.Run(test.path, func(t *testing.T) {
			// Action
			content, err := srv.ReadFile(context.Background(), test.path)

			// Assert
			asserts.Equals(t, test.expectedResponse, content)
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}

func TestService_buildFileName(t *testing.T) {
	testCases := []struct {
		filename         string
		content          []byte
		expectedFilename string
	}{
		{
			filename:         "test.txt",
			content:          []byte("test"),
			expectedFilename: "098f6bcd4621d373cade4e832627b4f6.txt",
		},
		{
			filename:         "test1.jpg",
			content:          []byte("test"),
			expectedFilename: "098f6bcd4621d373cade4e832627b4f6.jpg",
		},
		{
			filename:         "test2.jpg",
			content:          []byte("test2"),
			expectedFilename: "ad0234829205b9033196ba818f7a872b.jpg",
		},
		{
			filename:         "test2.txt",
			content:          nil,
			expectedFilename: "d41d8cd98f00b204e9800998ecf8427e.txt",
		},
	}

	for _, test := range testCases {
		t.Run(test.filename, func(t *testing.T) {
			// Arrange
			srv := New("./testdata")

			// Act
			actualFilename := srv.buildFileName(context.Background(), test.filename, test.content)

			// Assert
			asserts.Equals(t, test.expectedFilename, actualFilename)
		})
	}
}

func TestService_checkFileExists(t *testing.T) {
	testCases := []struct {
		path               string
		expectedFileExists bool
		expectedError      error
	}{
		{
			path:               "./testdata/test.txt",
			expectedFileExists: true,
			expectedError:      nil,
		},
		{
			path:               "./testdata/test1.txt",
			expectedFileExists: false,
			expectedError:      nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.path, func(t *testing.T) {
			// Arrange
			srv := New("./testdata")

			// Act
			fileExists, err := srv.checkFileExists(context.Background(), test.path)

			// Assert
			asserts.Equals(t, test.expectedFileExists, fileExists)
			asserts.Equals(t, test.expectedError, err)
		})
	}
}
