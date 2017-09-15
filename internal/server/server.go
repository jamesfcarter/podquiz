package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
	"github.com/jamesfcarter/podquiz/quiz"
)

type Server struct {
	Database *quiz.Database
	Template *assets.Templates
}

func (s *Server) App() (http.Handler, error) {
	mux := http.NewServeMux()
	assets.AddHandlers(mux)
	mux.HandleFunc("/guide.html", s.GuideHandler)
	mux.HandleFunc("/quiz.html", s.QuizHandler)
	mux.HandleFunc("/archive.html", s.ArchiveHandler)
	mux.HandleFunc("/comment", s.CommentHandler)
	mux.HandleFunc("/reload", s.ReloadHandler)
	mux.HandleFunc("/update", s.ReloadHandler)
	mux.HandleFunc("/quiz.php", s.QuizHandler)
	mux.HandleFunc("/podquiz.css", s.StylesheetHandler)
	mux.HandleFunc("/podquiz.rss", s.RSSHandler)
	mux.HandleFunc("/", s.IndexHandler)
	return mux, nil
}
