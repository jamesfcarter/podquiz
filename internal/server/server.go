package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jamesfcarter/podquiz/internal/assets"
	"github.com/jamesfcarter/podquiz/quiz"
)

type Server struct {
	Database   *quiz.Database
	Template   *assets.Templates
	PictureDir string
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("%s %s %s", r.Method, r.URL, elapsed)
	})
}

func (s *Server) App() (http.Handler, error) {
	mux := http.NewServeMux()
	assets.AddHandlers(mux)
	mux.HandleFunc("/guide.html", s.GuideHandler)
	mux.HandleFunc("/quiz.html", s.QuizHandler)
	mux.HandleFunc("/archive.html", s.ArchiveHandler)
	mux.HandleFunc("/archive.zip", s.ZipHandler)
	mux.HandleFunc("/comment", s.CommentHandler)
	mux.HandleFunc("/reload", s.ReloadHandler)
	mux.HandleFunc("/update", s.ReloadHandler)
	mux.HandleFunc("/quiz.php", s.QuizHandler)
	mux.HandleFunc("/podquiz.css", s.StylesheetHandler)
	mux.HandleFunc("/podquiz.rss", s.RSSHandler)
	mux.HandleFunc("/rss.php", s.RSSHandler)
	mux.HandleFunc("/pictures/", s.PictureHandler)
	mux.HandleFunc("/", s.IndexHandler)
	return Log(mux), nil
}
