package server

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/yosssi/gcss"
)

func (s *Server) RenderHTML(w http.ResponseWriter, template string, data interface{}) {
	s.render(w, "", template, data)
}

func (s *Server) RenderCSS(w http.ResponseWriter, template string, data interface{}) {
	s.render(w, "text/css", template, data)
}

func (s *Server) RenderSass(w http.ResponseWriter, template string, data interface{}) {
	var wg sync.WaitGroup
	rdr, wtr := io.Pipe()
	w.Header().Set("Content-Type", "text/css")
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer wtr.Close()
		err := s.Template.Exec(template, wtr, data)
		if err != nil {
			log.Print(err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := gcss.Compile(w, rdr)
		if err != nil {
			log.Print(err)
		}
	}()
	wg.Wait()
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
