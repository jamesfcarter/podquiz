package quiz_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

func TestFilename(t *testing.T) {
	quiz := &quiz.Episode{
		Number: 42,
	}
	cases := []struct {
		dir      string
		filename string
		comments string
	}{
		{"/foo/", "/foo/42.podcast", "/foo/42.comments"},
		{"/foo", "/foo/42.podcast", "/foo/42.comments"},
	}
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			fn := quiz.Filename(tc.dir)
			if fn != tc.filename {
				t.Errorf("expected %s, got %s", tc.filename, fn)
			}
			fn = quiz.CommentsFilename(tc.dir)
			if fn != tc.comments {
				t.Errorf("expected %s, got %s", tc.comments, fn)
			}
		})
	}
}

func TestRead(t *testing.T) {
	input := `42
PodQuiz 42
Mon, 5 Aug 1974 00:00:00 +0000
http://mp3.podquiz.com/pq42.mp3
12345,6789
Here is a description.
It has some lines.
`
	quiz, err := quiz.Read(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if quiz.Number != 42 {
		t.Errorf("bad quiz number %d", quiz.Number)
	}
	if quiz.Name != "PodQuiz 42" {
		t.Errorf("bad quiz number %s", quiz.Name)
	}
	loc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	if !quiz.Released.Equal(time.Date(1974, 8, 5, 0, 0, 0, 0, loc)) {
		t.Errorf("bad release time %v", quiz.Released)
	}
	if quiz.URL != "https://mp3.podquiz.com/pq42.mp3" {
		t.Errorf("bad quiz number %s", quiz.URL)
	}
	if quiz.Size != 12345 {
		t.Errorf("bad quiz size %d", quiz.Size)
	}
	if quiz.RestrictedSize != 6789 {
		t.Errorf("bad quiz restricted size %d", quiz.Size)
	}
	if !strings.Contains(quiz.Description, "some lines") {
		t.Errorf("bad descriptions %s", quiz.Description)
	}
}

func TestReadComments(t *testing.T) {
	input := `Mon, 5 Aug 1974 00:00:00 +0000
James
This is a test.
It has two lines
.
Mon, 5 Aug 1974 00:00:00 +0000
James2
This is another test.
.
`
	quiz := &quiz.Episode{}
	err := quiz.ReadComments(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(quiz.Comments) != 2 {
		t.Fatalf("expecting two comments, got %d", len(quiz.Comments))
	}
	loc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	if !quiz.Comments[0].Time.Equal(time.Date(1974, 8, 5, 0, 0, 0, 0, loc)) {
		t.Errorf("bad comment 1 time %v", quiz.Comments[0].Time)
	}
	if !quiz.Comments[1].Time.Equal(time.Date(1974, 8, 5, 0, 0, 0, 0, loc)) {
		t.Errorf("bad comment 2 time %v", quiz.Comments[1].Time)
	}
	if quiz.Comments[0].Author != "James" {
		t.Errorf("bad comment 1 author %s", quiz.Comments[0].Author)
	}
	if quiz.Comments[1].Author != "James2" {
		t.Errorf("bad comment 2 author %s", quiz.Comments[1].Author)
	}
	if quiz.Comments[0].Comment != "This is a test.\nIt has two lines\n" {
		t.Errorf("bad comment 1 %s", quiz.Comments[0].Comment)
	}
	if quiz.Comments[1].Comment != "This is another test.\n" {
		t.Errorf("bad comment 2 %s", quiz.Comments[1].Comment)
	}
}

func TestAddComment(t *testing.T) {
	dir, err := os.MkdirTemp("", "pqtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	loc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	fn := filepath.Join(dir, "42.comments")
	err = os.WriteFile(fn, []byte("hello\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	q := &quiz.Episode{
		Number: 42,
	}
	c := quiz.Comment{
		Time:    time.Date(1974, 8, 5, 0, 0, 0, 0, loc),
		Author:  "fred",
		Comment: "test\ncomment",
	}
	err = q.AddComment(dir, c)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.Comments) != 1 {
		t.Fatal("comment not added")
	}
	if q.Comments[0].Author != "fred" {
		t.Errorf("unexpected name %s", q.Comments[0].Author)
	}
	contents, err := os.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "hello\nMon, 5 Aug 1974 00:00:00 +0000\nfred\ntest\ncomment\n.\n" {
		t.Errorf("unexpected comments file contents: %s", string(contents))
	}
}
