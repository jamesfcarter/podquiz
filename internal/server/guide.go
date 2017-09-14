package server

import (
	"net/http"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) GuideHandler(w http.ResponseWriter, r *http.Request) {
	data := &assets.LayoutTemplateData{
		PageTitle: "Guide for Guest Hosts",
	}
	s.RenderHTML(w, "guide", data)
}
