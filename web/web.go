package web

import "embed"

//go:embed dist/*.* dist/css dist/js/*.js dist/pjs
var WebUI embed.FS
