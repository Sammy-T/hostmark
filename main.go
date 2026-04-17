package main

import (
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/sammy-t/hostmark/pwd"
	"gorm.io/gorm"
)

//go:embed template/readme.md
var readmeBytes []byte

var db *gorm.DB

var hmSecret string

func init() {
	const envPath = ".data/.env.local"
	var err error

	godotenv.Load(envPath)

	// Init the db
	cfg := &gorm.Config{
		TranslateError: true,
	}

	db, err = gorm.Open(sqlite.Open(".data/hostmark.db"), cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{}, &FailedLogin{}, &LockedToken{}, &RefreshToken{}, &Tag{}, &Note{}, &Preferences{})

	// Init env variable
	cwDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Current working directory: %v", err)
	}

	if hmSecret = os.Getenv("HM_SECRET"); hmSecret == "" {
		s := pwd.GenerateRandBytes(32)
		hmSecret = base64.StdEncoding.EncodeToString(s)

		entry := fmt.Sprintf("HM_SECRET=\"%v\"\n", hmSecret)

		ePath := filepath.Join(cwDir, envPath)
		if err = os.WriteFile(ePath, []byte(entry), 0644); err != nil {
			log.Fatalf("missing 'HM_SECRET': %v", err)
		}
	}

	// Init the readme file
	p := filepath.Join(cwDir, filesDir, "readme.md")
	if _, err = os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
		return
	}

	if err = os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		log.Printf("mkdir all: %v", err)
		return
	}

	if err = os.WriteFile(p, readmeBytes, 0644); err != nil {
		log.Printf("write file: %v", err)
	}
}

func main() {
	cwDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Current working directory: %v", err)
	}

	mux := http.NewServeMux()

	createFrontendHandlers(cwDir, mux)

	mux.HandleFunc("POST /api/auth/signup", handleSignup())
	mux.HandleFunc("POST /api/auth/login", handleLogin())
	mux.HandleFunc("GET /api/auth/refresh", handleRefresh())
	mux.HandleFunc("GET /api/auth/logout", handleLogout())

	mux.HandleFunc("GET /api/dir/{path...}", handleDirPath(cwDir))
	mux.HandleFunc("GET /api/file/{path...}", handleGetPath(cwDir))
	mux.HandleFunc("POST /api/file/{path...}", handlePostPath(cwDir))
	mux.HandleFunc("DELETE /api/file/{path...}", handleDelPath(cwDir))

	mux.HandleFunc("POST /api/note", handleCreateNote())
	mux.HandleFunc("GET /api/note/{id}", handleGetNote())
	mux.HandleFunc("POST /api/note/{id}", handleUpdateNote())
	mux.HandleFunc("DELETE /api/note/{id}", handleDelNote())
	mux.HandleFunc("GET /api/note/list", handleGetNotes())

	mux.HandleFunc("GET /api/tags", handleGetTags())

	mux.HandleFunc("POST /api/account", handleCreateUser())
	mux.HandleFunc("GET /api/account/me", handleGetMe())
	mux.HandleFunc("GET /api/account/{username}", handleGetUser())
	mux.HandleFunc("POST /api/account/{username}", handleUpdateUser())
	mux.HandleFunc("DELETE /api/account/{username}", handleDelUser())
	mux.HandleFunc("GET /api/account/list", handleGetUsers())

	addr := ":3000"

	log.Printf("Serving hostmark on %v", addr)
	log.Fatal(http.ListenAndServe(addr, logRequest(http.NewCrossOriginProtection().Handler(mux))))
}

// logRequest is middleware to log incoming request information.
//
// Note: Non-API routes are not logged.
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUrl := r.URL.String()

		if strings.Contains(reqUrl, "/api/") {
			log.Printf("%v %q %v from %v", r.Method, reqUrl, r.Proto, r.RemoteAddr)
		}

		next.ServeHTTP(w, r)
	})
}
