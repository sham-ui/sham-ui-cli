package assets

import (
	_ "embed"
)

//go:embed ssr.js
var Script []byte
