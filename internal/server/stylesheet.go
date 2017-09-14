package server

import (
	"net/http"
)

func (s *Server) StylesheetHandler(w http.ResponseWriter, r *http.Request) {
	s.RenderCSS(w, "stylesheet", nil)
}
