package server

import (
	"archive/zip"
	"context"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/context/ctxhttp"

	"github.com/jamesfcarter/podquiz/quiz"
)

func (s *Server) ZipHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	start := parseQuizFormNo(r, "start")
	end := parseQuizFormNo(r, "end")
	if start > end {
		x := start
		start = end
		end = x
	}
	w.Header().Set("Content-Type", "application/zip")
	zw := zip.NewWriter(w)
	ctx := r.Context()
	for n := start; n <= end; n++ {
		q := s.Database.Find(n)
		if q == nil {
			continue
		}
		err := zipQuiz(ctx, zw, q)
		if err != nil {
			log.Printf("failed to zip %s: %v", q.MP3(), err)
		}
	}
	err := zw.Close()
	if err != nil {
		log.Print(err)
	}
}

func zipQuiz(ctx context.Context, zw *zip.Writer, q *quiz.Episode) error {
	w, err := zw.Create(q.MP3())
	if err != nil {
		return err
	}
	r, err := ctxhttp.Get(ctx, nil, q.URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = io.Copy(w, r.Body)
	if err != nil {
		return err
	}
	return nil
}
