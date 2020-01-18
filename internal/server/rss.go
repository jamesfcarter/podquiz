package server

import (
	"net/http"
	"time"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	data := &assets.RSSTemplateData{
		LastBuild: s.Database.Find(s.Database.MostRecent()).Released,
		ThisYear:  time.Now().Year(),
		Quizzes:   s.Database.Page(s.Database.MostRecent(), s.Database.Count(760)),
	}
	s.RenderRSS(w, "rss", data)
}
