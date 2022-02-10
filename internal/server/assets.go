package server

import (
	"embed"
)

//go:embed public_html
var assets embed.FS
