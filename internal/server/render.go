package server

import (
	"fmt"
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

func (s *Server) RenderRSS(w http.ResponseWriter, template string, data interface{}) {
	s.render(w, "application/rss+xml", template, data)
}

func (s *Server) RenderSass(w http.ResponseWriter, template string, data interface{}) {
	var wg sync.WaitGroup
	rdr, wtr := io.Pipe()
	w.Header().Set("Content-Type", "text/css")
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer wtr.Close()
		err := s.executeTemplate(wtr, template, data)
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

func (s *Server) executeTemplate(w io.Writer, name string, data interface{}) error {
	template, ok := s.Templates[name]
	if !ok {
		return fmt.Errorf("no such template %s", name)
	}
	return template.Execute(w, data)
}

func (s *Server) render(w http.ResponseWriter, contentType string, template string, data interface{}) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	err := s.executeTemplate(w, template, data)
	if err != nil {
		log.Print(err)
	}
}
