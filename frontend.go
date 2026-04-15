//go:build !dev

package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	httpExt "github.com/sammy-t/hostmark/internal/http"
)

func createFrontendHandlers(cwDir string, mux *http.ServeMux) {
	buildSite("pnpm", cwDir)

	fs := httpExt.FileSys{FileSystem: http.Dir(filepath.Join(cwDir, "/web/build"))}

	mux.Handle("/", http.FileServer(fs))
	mux.HandleFunc("/file/{path...}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(cwDir, "/web/build/file.html"))
	})
	mux.HandleFunc("/note/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(cwDir, "/web/build/note/id.html"))
	})
}

func buildSite(pkgManager string, cwDir string) {
	cmdDir := filepath.Join(cwDir, "web")

	cmdBuild := exec.Command(pkgManager, "run", "build")
	cmdBuild.Dir = cmdDir
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr

	// Run build and await
	if err := cmdBuild.Run(); err != nil {
		if pkgManager != "npm" {
			log.Printf("%v build failed: %v.\nFalling back to npm...", pkgManager, err)

			buildSite("npm", cwDir)
		} else {
			log.Fatalf("build frontend: %v", err)
		}
	}
}
