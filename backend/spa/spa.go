package spa

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed frontend/build/*
var embeddedFiles embed.FS

// Handler serves the embedded SPA.
func Handler() http.HandlerFunc {
	distFS, err := fs.Sub(embeddedFiles, "frontend/build")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		f, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			r.URL.Path = "/"
		} else {
			f.Close()
		}

		if strings.HasPrefix(r.URL.Path, "/_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Cache-Control", "no-cache")
		}

		fileServer.ServeHTTP(w, r)
	}
}
