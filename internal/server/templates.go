package server

import (
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

type LayoutTemplateData struct {
	PageTitle string
}

type ArchiveTemplateData struct {
	PageTitle  string
	MostRecent int
	Quizzes    []*quiz.Episode
}

type IndexTemplateData struct {
	PageTitle string
	Quizzes   []*quiz.Episode
}

type QuizTemplateData struct {
	PageTitle   string
	CommentName string
	Comment     string
	Quiz        *quiz.Episode
}

type DoneTemplateData struct {
	PageTitle string
	Done      map[string]string
}

type RSSTemplateData struct {
	LastBuild  time.Time
	ThisYear   int
	Quizzes    []*quiz.Episode
	Restricted bool
}
