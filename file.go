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
	"strings"

	"github.com/sammy-t/hostmark/internal/auth"
	"gorm.io/gorm"
)

const filesDir string = ".files"

type PathEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func handleDirPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResFile, auth.PermRead, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResFile, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		urlPath := r.PathValue("path")

		entries, err := os.ReadDir(filepath.Join(cwDir, filesDir, urlPath))
		if err != nil {
			log.Printf("read dir: %v", err)
			http.Error(w, "unable to read directory", http.StatusInternalServerError)
			return
		}

		var pathEntries []PathEntry

		for _, entry := range entries {
			if containsDotFile(entry.Name()) {
				continue
			}

			p := PathEntry{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
			}

			pathEntries = append(pathEntries, p)
		}

		resp, err := json.Marshal(pathEntries)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleGetPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResFile, auth.PermRead, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResFile, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		urlPath := r.PathValue("path")

		if containsDotFile(urlPath) {
			log.Printf("requested ignored path %q", urlPath)
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		p := filepath.Join(cwDir, filesDir, urlPath)

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
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResFile, auth.PermUpdate, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResFile, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		urlPath := r.PathValue("path")

		if containsDotFile(urlPath) {
			log.Printf("requested ignored path %q", urlPath)
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()

		if err != nil {
			log.Printf("read req body: %v", err)
			http.Error(w, "unable to read body", http.StatusInternalServerError)
			return
		}

		p := filepath.Join(cwDir, filesDir, urlPath)

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
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResFile, auth.PermDelete, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResFile, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		urlPath := r.PathValue("path")

		if containsDotFile(urlPath) {
			log.Printf("requested ignored path %q", urlPath)
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		if err := os.Remove(filepath.Join(cwDir, filesDir, urlPath)); err != nil {
			log.Printf("remove file: %v", err)
			http.Error(w, "unable to remove file", http.StatusInternalServerError)
			return
		}
	}
}

// containsDotFile is a helper to determine if the names
// of the file or its parent directories contain a leading dot.
func containsDotFile(pathName string) bool {
	for part := range strings.SplitSeq(pathName, "/") {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}

	return false
}
