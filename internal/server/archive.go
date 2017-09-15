package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	data := &assets.IndexTemplateData{
		PageTitle: "Podquiz Archive",
		Quizzes:   s.Database.All(),
	}
	s.RenderHTML(w, "archive", data)
}
