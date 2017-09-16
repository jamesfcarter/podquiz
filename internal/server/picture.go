package server

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) PictureHandler(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Base(r.URL.Path)
	if !strings.HasSuffix(fn, ".jpg") && !strings.HasSuffix(fn, ".png") {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	path := filepath.Join(s.PictureDir, fn)
	http.ServeFile(w, r, path)
}
