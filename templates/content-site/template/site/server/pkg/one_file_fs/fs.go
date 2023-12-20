package one_file_fs

import (
	"net/http"
)

type fileSystem struct {
	file *file
}

func (fs fileSystem) Open(_ string) (http.File, error) {
	return fs.file, nil
}

func New(path string, contents []byte) *fileSystem {
	return &fileSystem{
		file: &file{
			stats: stats{
				name: path,
				size: int64(len(contents)),
			},
			contents: contents,
			offset:   0,
		},
	}
}
