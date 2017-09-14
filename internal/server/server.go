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
	mux.HandleFunc("/podquiz.css", s.StylesheetHandler)
	return mux, nil
}
