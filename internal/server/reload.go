package server

import (
	"fmt"
	"net/http"
)

func (s *Server) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	then := s.Database.MostRecent()
	s.Database.Update()
	now := s.Database.MostRecent()
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Moved from %d to %d", then, now)
}
