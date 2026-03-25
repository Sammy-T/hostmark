package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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

	if dev {
		startDevServer("pnpm", cwDir)
		http.Handle("/", createDevHandler())
	} else {
		buildSite("pnpm")

		fs := httpExt.FileSys{FileSystem: http.Dir(filepath.Join(cwDir, "/web/build"))}

		http.Handle("/", http.FileServer(fs))
		http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(cwDir, "/web/build/login.html"))
		})
	}

	http.HandleFunc("POST /api/auth/signup", handleSignUp())
	http.HandleFunc("POST /api/auth/login", handleLogIn())

	http.HandleFunc("GET /api/dir/{path...}", handleDirPath(cwDir))
	http.HandleFunc("GET /api/file/{path...}", handleGetPath(cwDir))
	http.HandleFunc("POST /api/file/{path...}", handlePostPath(cwDir))
	http.HandleFunc("DELETE /api/file/{path...}", handleDelPath(cwDir))

	addr := ":3000"

	log.Printf("Serving hostmark on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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
