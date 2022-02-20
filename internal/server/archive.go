package server

import (
	"net/http"
)

func (s *Server) ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	data := &ArchiveTemplateData{
		PageTitle:  "Podquiz Archive",
		MostRecent: s.Database.MostRecent(),
		Quizzes:    s.Database.All(),
	}
	s.RenderHTML(w, "archive.html", data)
}
