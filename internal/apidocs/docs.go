package apidocs

import (
	"embed"
	"net/http"
)

//go:embed index.html openapi.json swagger-initializer.js
var files embed.FS

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /docs", serve("index.html", "text/html; charset=utf-8"))
	mux.HandleFunc("GET /docs/", serve("index.html", "text/html; charset=utf-8"))
	mux.HandleFunc(
		"GET /docs/swagger-initializer.js",
		serve("swagger-initializer.js", "application/javascript; charset=utf-8"),
	)
	mux.HandleFunc("GET /openapi.json", serve("openapi.json", "application/json"))
}

func serve(name string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		content, err := files.ReadFile(name)
		if err != nil {
			http.Error(w, "documentation is unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	}
}
