package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var dev bool

type fileSys struct {
	http.FileSystem
}

type PathEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
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

		fs := fileSys{http.Dir(filepath.Join(cwDir, "/web/build"))}

		http.Handle("/", http.FileServer(fs))
		http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(cwDir, "/web/build/login.html"))
		})
	}

	http.HandleFunc("GET /dir/{path...}", handleDirPath(cwDir))

	addr := ":3000"

	log.Printf("Serving hostmark on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleDirPath(cwDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirPath := r.PathValue("path")

		log.Printf("req %q %q", r.URL.String(), dirPath)

		entries, err := os.ReadDir(filepath.Join(cwDir, ".files", dirPath))
		if err != nil {
			log.Fatal(err)
		}

		var pathEntries []PathEntry

		for _, entry := range entries {
			p := PathEntry{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
			}

			pathEntries = append(pathEntries, p)
		}

		log.Print(pathEntries)

		jsonBytes, err := json.Marshal(pathEntries)

		w.Write(jsonBytes)
	}
}

func (fs fileSys) Open(name string) (http.File, error) {
	// log.Printf("file %q", name)

	re := regexp.MustCompile(`\/[\w\-]+$`)

	// Attempt to open an html file matching the directory.
	// i.e. /hello => /hello.html
	//
	// By default, Go's FileServer supports directories containing index.html files.
	// See: https://pkg.go.dev/net/http@go1.26.0#FileServer
	if re.MatchString(name) {
		htmlName := name + ".html"
		log.Printf("%q => %q", name, htmlName)

		file, err := fs.FileSystem.Open(htmlName)
		if err != nil {
			log.Printf("file %q: %v", htmlName, err)
		} else {
			return file, nil
		}
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		log.Printf("file %q: %v", name, err)
		return nil, err
	}

	return file, nil
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

func startDevServer(pkgManager string, cwDir string) {
	log.Println("Starting Vite dev server...")

	cmdDir := filepath.Join(cwDir, "web")

	cmdInst := exec.Command(pkgManager, "i")
	cmdInst.Dir = cmdDir
	cmdInst.Stdout = os.Stdout
	cmdInst.Stderr = os.Stderr

	var err error

	// Run install and await
	if err = cmdInst.Run(); err != nil {
		if pkgManager != "npm" {
			log.Printf("%v install failed: %v.\nFalling back to npm...", pkgManager, err)

			startDevServer("npm", cwDir)
			return
		}

		log.Fatalf("%v install: %v", pkgManager, err)
	}

	cmd := exec.Command(pkgManager, "run", "dev")
	cmd.Dir = cmdDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run dev without awaiting
	if err = cmd.Start(); err != nil {
		log.Fatalf("Dev server: %v", err)
	}

	log.Printf("Dev server running as pid %v", cmd.Process.Pid)
}

func createDevHandler() http.Handler {
	devUrl, err := url.Parse("http://localhost:5173")
	if err != nil {
		log.Fatalf("Dev url: %v", err)
	}

	errHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Print(err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadGateway)

		w.Write([]byte("Error: " + err.Error()))
	}

	rewrite := func(r *httputil.ProxyRequest) {
		// log.Printf("proxy req path %q", r.In.URL.String())

		if r.In.URL.Path == "/auth" {
			// Can't directly set r.Out.URL here for some reason
			// even though it's a pointer
			r.Out.URL.Scheme = devUrl.Scheme
			r.Out.URL.Host = devUrl.Host
			r.Out.URL.Path = "/login"

			log.Printf("%q => %q %q", r.In.URL.String(), r.Out.URL.String(), r.Out.URL.Path)
			return
		}

		r.SetURL(devUrl)
		// log.Printf("=> %q", r.Out.URL.String())
	}

	proxy := &httputil.ReverseProxy{
		Rewrite:      rewrite,
		ErrorHandler: errHandler,
	}

	return proxy
}
