package assets

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mattetti/filebuffer"
)

func AddHandlers(mux *http.ServeMux) {
	for _, name := range AssetNames() {
		if !staticAsset(name) {
			continue
		}
		log.Printf("%s: %s", name, webPath(name))
		mux.HandleFunc(webPath(name), AssetWriterFunc(name))
	}
}

func AssetWriterFunc(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := filebuffer.New(MustAsset(name))
		info, _ := AssetInfo(name)
		http.ServeContent(w, r, webPath(name), info.ModTime(), buf)
	}
}

func staticAsset(name string) bool {
	parts := strings.SplitN(name, "/", 2)
	return parts[0] == "static"
}

func webPath(name string) string {
	parts := strings.Split(name, "/")
	return "/" + filepath.Join(parts[1:]...)
}
