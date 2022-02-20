package assets

import (
	"embed"
	"html/template"
	"io/fs"
	"strings"
	"time"
)

//go:embed views/*
var views embed.FS

func Templates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	files, err := fs.ReadDir(views, "views")
	if err != nil {
		return nil, err
	}
	for _, tmpl := range files {
		name := tmpl.Name()
		if strings.HasSuffix(name, ".tmpl") {
			continue
		}
		parsed, err := template.New("").
			Funcs(templateFuncs).
			ParseFS(views, "views/"+name, "views/*.tmpl")
		if err != nil {
			return nil, err
		}

		templates[name] = parsed.Lookup(name)
	}
	return templates, nil
}

var templateFuncs template.FuncMap = template.FuncMap{
	"rssTime":   timeFormat(rssTimeFormat),
	"indexTime": timeFormat(indexTimeFormat),
	"xml":       xml,
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

func xml() template.HTML {
	return template.HTML(`<?xml version="1.0"?>`)
}
