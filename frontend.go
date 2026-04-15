//go:build !dev

package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	fsExt "github.com/sammy-t/hostmark/internal/fs"
)

//go:embed web/build/*
var content embed.FS

func createFrontendHandlers(_ string, mux *http.ServeMux) {
	fsys, err := fs.Sub(content, "web/build")
	if err != nil {
		log.Fatalf("frontend: %v", err)
	}

	frontend := fsExt.FS{FS: fsys}

	mux.Handle("/", http.FileServerFS(frontend))
	mux.HandleFunc("/file/{path...}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend, "file.html")
	})
	mux.HandleFunc("/note/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, frontend, "note/id.html")
	})
}
