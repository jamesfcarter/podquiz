package server

import (
	"net/http"
)

func (s *Server) DoneHandler(w http.ResponseWriter, r *http.Request) {
	data := &DoneTemplateData{
		PageTitle: "Rounds Already Done",
		Done:      s.Done.Rounds,
	}
	s.RenderHTML(w, "done.html", data)
}
