package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := &assets.IndexTemplateData{
		PageTitle: "Podquiz",
		Quizzes:   s.Database.Page(s.Database.MostRecent(), 10),
	}
	s.RenderHTML(w, "index", data)
}
