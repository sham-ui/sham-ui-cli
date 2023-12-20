package one_file_fs

import "errors"

var (
	errReadDir      = errors.New("can't read dir")
	errSeekNegative = errors.New("can't seek negative position")
	errSeekInvalid  = errors.New("can't seek invalid position")
)
