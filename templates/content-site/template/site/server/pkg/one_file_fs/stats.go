package one_file_fs

import (
	"io/fs"
	"time"
)

type stats struct {
	name string
	size int64
}

func (f stats) Name() string       { return f.name }
func (f stats) Size() int64        { return f.size }
func (f stats) Mode() fs.FileMode  { return fs.ModePerm }
func (f stats) ModTime() time.Time { return time.Time{} }
func (f stats) IsDir() bool        { return false }
func (f stats) Sys() any           { return nil }
