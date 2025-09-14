package server

import (
	"net/http"
	"time"
)

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	data := &RSSTemplateData{
		LastBuild: s.Database.Find(s.Database.MostRecent()).Released,
		ThisYear:  time.Now().Year(),
		Quizzes:   s.Database.Page(s.Database.MostRecent(), s.Database.Count(760)),
	}
	s.RenderRSS(w, "podquiz.rss", data)
}

func (s *Server) RSSFullHandler(w http.ResponseWriter, r *http.Request) {
	data := &RSSTemplateData{
		LastBuild: s.Database.Find(s.Database.MostRecent()).Released,
		ThisYear:  time.Now().Year(),
		Quizzes:   s.Database.All(),
	}
	s.RenderRSS(w, "podquiz.rss", data)
}

func (s *Server) RSSRestrictedHandler(w http.ResponseWriter, r *http.Request) {
	data := &RSSTemplateData{
		LastBuild:  s.Database.Find(s.Database.MostRecent()).Released,
		ThisYear:   time.Now().Year(),
		Quizzes:    s.Database.Page(s.Database.MostRecent(), s.Database.Count(1061)),
		Restricted: true,
	}
	s.RenderRSS(w, "podquiz.rss", data)
}
