package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) QuizHandler(w http.ResponseWriter, r *http.Request) {
	q := s.Database.Find(parseQuizNo(r.URL.Query()["q"]))
	if q == nil {
		return
	}
	data := &assets.QuizTemplateData{
		PageTitle: q.Name,
		Quiz:      q,
	}
	s.RenderHTML(w, "quiz", data)
}

func parseQuizNo(arg []string) int {
	parts := strings.Split(arg[0], "/")
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}
