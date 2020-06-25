package server

import (
	"fmt"
	"net/http"
)

func (s *Server) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := s.Done.Load()
	if err != nil {
		fmt.Fprintf(w, "Error loading donefile: %v\n", err)
	}
	then := s.Database.MostRecent()
	s.Database.Update()
	now := s.Database.MostRecent()
	fmt.Fprintf(w, "Moved from %d to %d", then, now)
}
