package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) DoneHandler(w http.ResponseWriter, r *http.Request) {
	data := &assets.DoneTemplateData{
		PageTitle: "Rounds Already Done",
		Done:      s.Done.Rounds,
	}
	s.RenderHTML(w, "done", data)
}
