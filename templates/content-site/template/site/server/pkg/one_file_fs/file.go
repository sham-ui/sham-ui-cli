package one_file_fs

import (
	"io"
	"io/fs"
)

type file struct {
	stats    stats
	contents []byte
	offset   int64
}

func (o *file) Close() error { return nil }

func (o *file) Read(p []byte) (int, error) {
	if o.offset >= int64(len(o.contents)) {
		return 0, io.EOF
	}
	n := copy(p, o.contents[o.offset:])
	o.offset += int64(n)
	return n, nil
}

func (o *file) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = o.offset + offset
	case io.SeekEnd:
		abs = int64(len(o.contents)) + offset
	default:
		return 0, errSeekInvalid
	}
	if abs < 0 {
		return 0, errSeekNegative
	}
	o.offset = abs
	return abs, nil
}

func (o *file) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, errReadDir
}

func (o file) Stat() (fs.FileInfo, error) {
	return o.stats, nil
}
