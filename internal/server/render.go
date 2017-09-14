package server

import (
	"log"
	"net/http"
)

func (s *Server) RenderHTML(w http.ResponseWriter, template string, data interface{}) {
	s.render(w, "", template, data)
}

func (s *Server) RenderCSS(w http.ResponseWriter, template string, data interface{}) {
	s.render(w, "text/css", template, data)
}

func (s *Server) render(w http.ResponseWriter, contentType string, template string, data interface{}) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	err := s.Template.Exec(template, w, data)
	if err != nil {
		log.Print(err)
	}
}
