package server

import (
	"net/http"
)

func (s *Server) GuideHandler(w http.ResponseWriter, r *http.Request) {
	data := &LayoutTemplateData{
		PageTitle: "Guide for Guest Hosts",
	}
	s.RenderHTML(w, "guide.html", data)
}
