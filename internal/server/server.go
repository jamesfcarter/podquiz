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
	return mux, nil
}
