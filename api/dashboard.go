package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

func dashboardHandler(directory string) http.Handler {
	files := http.FileServer(http.Dir(directory))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "GET" || r.Method == "HEAD") && (r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/factories/") || strings.HasPrefix(r.URL.Path, "/stories/")) {
			http.ServeFile(w, r, filepath.Join(directory, "index.html"))
			return
		}
		files.ServeHTTP(w, r)
	})
}
