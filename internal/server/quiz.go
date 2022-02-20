package server

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) QuizHandler(w http.ResponseWriter, r *http.Request) {
	q := s.Database.Find(parseQuizNo(r.URL.Query()["q"]))
	if q == nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	data := &QuizTemplateData{
		PageTitle: q.Name,
		Quiz:      q,
	}
	s.RenderHTML(w, "quiz.html", data)
}

func parseQuizNo(arg []string) int {
	if len(arg) == 0 {
		return 0
	}
	parts := strings.Split(arg[0], "/")
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}
