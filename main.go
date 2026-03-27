package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	httpExt "github.com/sammy-t/hostmark/internal/http"
	"gorm.io/gorm"
)

var dev bool
var db *gorm.DB

func init() {
	var err error

	cfg := &gorm.Config{
		TranslateError: true,
	}

	db, err = gorm.Open(sqlite.Open("hostmark.db"), cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{}, &FailedLogin{}, &LockedToken{}, &RefreshToken{})
}

func main() {
	flag.BoolVar(&dev, "dev", false, "development mode")
	flag.Parse()

	cwDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Current working directory: %v", err)
	}

	mux := http.NewServeMux()

	if dev {
		startDevServer("pnpm", cwDir)
		mux.Handle("/", createDevHandler())
	} else {
		buildSite("pnpm")

		fs := httpExt.FileSys{FileSystem: http.Dir(filepath.Join(cwDir, "/web/build"))}

		mux.Handle("/", http.FileServer(fs))
		mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(cwDir, "/web/build/login.html"))
		})
	}

	mux.HandleFunc("POST /api/auth/signup", handleSignup())
	mux.HandleFunc("POST /api/auth/login", handleLogin())
	mux.HandleFunc("GET /api/auth/refresh", handleRefresh())
	mux.HandleFunc("GET /api/auth/logout", handleLogout())

	mux.HandleFunc("GET /api/dir/{path...}", handleDirPath(cwDir))
	mux.HandleFunc("GET /api/file/{path...}", handleGetPath(cwDir))
	mux.HandleFunc("POST /api/file/{path...}", handlePostPath(cwDir))
	mux.HandleFunc("DELETE /api/file/{path...}", handleDelPath(cwDir))

	mux.HandleFunc("GET /api/account/me", handleGetMe())

	addr := ":3000"

	log.Printf("Serving hostmark on %v", addr)
	log.Fatal(http.ListenAndServe(addr, logRequest(mux)))
}

func buildSite(pkgManager string) {
	cwDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Current working directory: %v", err)
	}

	cmdDir := filepath.Join(cwDir, "web")

	cmdBuild := exec.Command(pkgManager, "run", "build")
	cmdBuild.Dir = cmdDir
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr

	// Run build and await
	if err = cmdBuild.Run(); err != nil {
		if pkgManager != "npm" {
			log.Printf("%v build failed: %v.\nFalling back to npm...", pkgManager, err)

			buildSite("npm")
		}
	}
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
