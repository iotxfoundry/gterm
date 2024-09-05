package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed dist/*.* dist/css dist/js/*.js dist/pjs
var WebUI embed.FS

func WebServer() (h http.Handler) {
	var hf http.FileSystem
	_, err := os.Stat("web/dist")
	if err != nil {
		fsys, err := fs.Sub(WebUI, "dist")
		if err != nil {
			return nil
		}
		hf = http.FS(fsys)
	} else {
		hf = http.Dir("web/dist")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		http.FileServer(hf).ServeHTTP(w, r)
	})
}
