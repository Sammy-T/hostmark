package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PathEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func handleDirPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.PathValue("path")

		entries, err := os.ReadDir(filepath.Join(cwDir, ".files", urlPath))
		if err != nil {
			log.Printf("read dir: %v", err)
			http.Error(w, "unable to read directory", http.StatusInternalServerError)
			return
		}

		var pathEntries []PathEntry

		for _, entry := range entries {
			p := PathEntry{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
			}

			pathEntries = append(pathEntries, p)
		}

		jsonBytes, err := json.Marshal(pathEntries)
		w.Write(jsonBytes)
	}
}

func handleGetPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.PathValue("path")

		p := filepath.Join(cwDir, ".files", urlPath)

		info, err := os.Stat(p)
		if err != nil {
			msg := "unable to read file"
			code := http.StatusInternalServerError

			if errors.Is(err, fs.ErrNotExist) {
				msg = "file not found"
				code = http.StatusNotFound
			}

			log.Printf("stat file: %v", err)
			http.Error(w, msg, code)
			return
		}

		if info.IsDir() {
			msg := "requested file is a directory"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := os.ReadFile(p)
		if err != nil {
			log.Printf("read file: %v", err)
			http.Error(w, "unable to read file", http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

func handlePostPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.PathValue("path")

		body, err := io.ReadAll(r.Body)
		r.Body.Close()

		if err != nil {
			log.Printf("read req body: %v", err)
			http.Error(w, "unable to read body", http.StatusInternalServerError)
			return
		}

		p := filepath.Join(cwDir, ".files", urlPath)

		if err = os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			log.Printf("mkdir all: %v", err)
			http.Error(w, "unable to create directory", http.StatusInternalServerError)
			return
		}

		if err = os.WriteFile(p, body, 0644); err != nil {
			log.Printf("write file: %v", err)
			http.Error(w, "unable to write file", http.StatusInternalServerError)
			return
		}
	}
}

func handleDelPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.PathValue("path")

		if err := os.Remove(filepath.Join(cwDir, ".files", urlPath)); err != nil {
			log.Printf("remove file: %v", err)
			http.Error(w, "unable to remove file", http.StatusInternalServerError)
			return
		}
	}
}
