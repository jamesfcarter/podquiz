package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if !validIndexPath(r.URL.Path) {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	data := &assets.IndexTemplateData{
		PageTitle: "Podquiz",
		Quizzes:   s.Database.Page(s.Database.MostRecent(), 10),
	}
	s.RenderHTML(w, "index", data)
}

func validIndexPath(path string) bool {
	for _, p := range []string{
		"/",
		"/index.html",
		"/index.php",
	} {
		if path == p {
			return true
		}
	}
	return false
}
