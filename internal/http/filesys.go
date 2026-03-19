package http

import (
	"log"
	"net/http"
	"regexp"
)

// FileSys embeds http.FileSystem to extend its functionality.
type FileSys struct {
	http.FileSystem
}

// Open extends FileSys's http.FileSystem
// to include `/<dir-name>.html` matching
// in addition to the default `/<dir-name>/index.html` matching.
func (fs FileSys) Open(name string) (http.File, error) {
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
