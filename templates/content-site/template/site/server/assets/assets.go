package assets

import (
	"embed"
	"io/fs"
)

//go:embed files/*
var files embed.FS

func Files() (fs.FS, error) {
	return fs.Sub(files, "files")
}
