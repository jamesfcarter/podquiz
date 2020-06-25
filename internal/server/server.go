package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jamesfcarter/podquiz/internal/assets"
	"github.com/jamesfcarter/podquiz/internal/done"
	"github.com/jamesfcarter/podquiz/quiz"
)

type Server struct {
	Database   *quiz.Database
	Template   *assets.Templates
	PictureDir string
	MerchUrl   string
	Done       *done.Done
}

type loggedResponseWriter struct {
	http.ResponseWriter
	Code int
}

func (l *loggedResponseWriter) WriteHeader(code int) {
	l.Code = code
	l.ResponseWriter.WriteHeader(code)
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggedResponseWriter{
			ResponseWriter: w,
			Code:           http.StatusOK,
		}
		start := time.Now()
		handler.ServeHTTP(lw, r)
		elapsed := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL, lw.Code, elapsed)
	})
}

func (s *Server) App() (http.Handler, error) {
	mux := http.NewServeMux()
	assets.AddHandlers(mux)
	merchRedirect := NewRedirectHandler(s.MerchUrl)
	mux.HandleFunc("/merch/", merchRedirect)
	mux.HandleFunc("/guide.html", s.GuideHandler)
	mux.HandleFunc("/done.html", s.DoneHandler)
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
	mux.HandleFunc("/full.rss", s.RSSFullHandler)
	mux.HandleFunc("/pictures/", s.PictureHandler)
	mux.HandleFunc("/", s.IndexHandler)
	return Log(mux), nil
}
