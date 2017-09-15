package assets

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

type LayoutTemplateData struct {
	PageTitle string
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

type RSSTemplateData struct {
	LastBuild time.Time
	ThisYear  int
	Quizzes   []*quiz.Episode
}

// Templates is a simple map of named templates".
type Templates map[string]*template.Template

// MakeTemplates constructs a Templates object from asset data.
func MakeTemplates() (*Templates, error) {
	result := make(Templates)
	for name, tmpl := range map[string]string{
		"index":      withLayout("index"),
		"quiz":       withLayout("quiz"),
		"guide":      withLayout("guide"),
		"stylesheet": stringAsset("stylesheet"),
		"rss":        stringAsset("rss"),
	} {
		t, err := template.New(name).Funcs(templateFuncs).Parse(tmpl)
		if err != nil {
			return nil, err
		}
		result[name] = t
	}
	return &result, nil
}

// Execute executes the named template.
func (t *Templates) Exec(name string, w io.Writer, data interface{}) error {
	tmpl, ok := (*t)[name]
	if !ok {
		return fmt.Errorf("no such template %s", name)
	}
	return tmpl.Execute(w, data)
}

func withLayout(tmpl string) string {
	return stringAsset("layout") + stringAsset(tmpl)
}

func stringAsset(name string) string {
	return string(MustAsset(assetName(name)))
}

func assetName(tmpl string) string {
	return "views/" + tmpl + ".tmpl"
}

var templateFuncs template.FuncMap = template.FuncMap{
	"rssTime":   timeFormat(rssTimeFormat),
	"indexTime": timeFormat(indexTimeFormat),
}

const (
	rssTimeFormat   = "Mon, 02 Jan 2006 15:04:05 -0700"
	indexTimeFormat = "2006-01-02 15:04"
)

func timeFormat(fmt string) func(t time.Time) template.HTML {
	return func(t time.Time) template.HTML {
		return template.HTML(t.Format(fmt))
	}
}
