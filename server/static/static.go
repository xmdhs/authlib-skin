package static

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed files
var staticFs embed.FS

func StaticServer() http.Handler {
	serverRoot, err := fs.Sub(staticFs, "files")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Get("/", index)
	r.Get("/*", index)

	r.Mount("/assets", http.FileServer(http.FS(serverRoot)))

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	b, err := staticFs.ReadFile("files/index.html")
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
