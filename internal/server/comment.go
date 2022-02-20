package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

func (s *Server) CommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q := s.Database.Find(parseQuizFormNo(r, "q"))
	if q == nil {
		http.Error(w, "bad quiz", http.StatusBadRequest)
		return
	}
	data := &QuizTemplateData{
		PageTitle: q.Name,
		Quiz:      q,
	}
	from := r.Form.Get("name")
	comment := r.Form.Get("comment")
	if from != "" && comment != "" && parseQuizFormNo(r, "human") == q.Number {
		q.AddComment(s.Database.Dir, quiz.Comment{
			Time:    time.Now(),
			Author:  from,
			Comment: comment,
		})
	} else {
		data.CommentName = from
		data.Comment = comment
	}
	s.RenderHTML(w, "quiz.html", data)
}

func parseQuizFormNo(r *http.Request, arg string) int {
	parts := strings.Split(r.Form.Get(arg), "/")
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}
