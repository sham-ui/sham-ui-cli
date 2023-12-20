package one_file_fs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"site/pkg/asserts"
)

// Compile time check that file implements http.File.
var _ http.File = (*file)(nil)

// Compile time check that fileSystem implements http.FileSystem.
var _ http.FileSystem = (*fileSystem)(nil)

func TestFileSystemServe(t *testing.T) {
	content := []byte("1234567890")
	path := "content/test.txt"

	fs := New(path, content)
	req := httptest.NewRequest(http.MethodGet, "/file", nil)
	resp := httptest.NewRecorder()
	http.FileServer(fs).ServeHTTP(resp, req)

	asserts.Equals(t, "1234567890", resp.Body.String(), "content")
	asserts.Equals(t, http.StatusOK, resp.Code, "status code")
}
