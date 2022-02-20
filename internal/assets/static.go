package assets

import (
	"bytes"
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//go:embed static
var static embed.FS

func AddStaticHandlers(mux *http.ServeMux) {
	if err := fs.WalkDir(static, "static", addAsset(mux)); err != nil {
		log.Fatalf("failed to read static assets: %v", err)
	}
}

func addAsset(mux *http.ServeMux) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if path == "static" {
			return nil
		}
		urlPath := strings.TrimPrefix(path, "static")
		mux.HandleFunc(urlPath, staticServe(path))
		return nil
	}
}

func staticServe(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := static.Open(path)
		if err != nil {
			log.Printf("unable to open static asset %s: %v", path, err)
		}
		fInfo, err := f.Stat()
		if err != nil {
			log.Printf("unable to stat static asset %s: %v", path, err)
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Printf("unable to read static asset %s: %v", path, err)
		}
		http.ServeContent(w, r, path, fInfo.ModTime(), bytes.NewReader(buf))
	}
}
