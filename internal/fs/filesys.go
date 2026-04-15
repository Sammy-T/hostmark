package fs

import (
	"io/fs"
	"log"
	"regexp"
)

// FS embeds fs.FS to extend its functionality.
type FS struct {
	fs.FS
}

// Open extends FS's fs.FS
// to include `/<dir-name>.html` matching
// in addition to the default `/<dir-name>/index.html` matching.
func (fs FS) Open(name string) (fs.File, error) {
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

		file, err := fs.FS.Open(htmlName)
		if err != nil {
			log.Printf("file %q: %v", htmlName, err)
		} else {
			return file, nil
		}
	}

	file, err := fs.FS.Open(name)
	if err != nil {
		log.Printf("file %q: %v", name, err)
		return nil, err
	}

	return file, nil
}
