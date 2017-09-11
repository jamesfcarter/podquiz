package quiz_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

func TestFilename(t *testing.T) {
	quiz := &quiz.Quiz{
		Number: 42,
	}
	cases := []struct {
		dir      string
		expected string
	}{
		{"/foo/", "/foo/42.podcast"},
		{"/foo", "/foo/42.podcast"},
	}
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			fn := quiz.Filename(tc.dir)
			if fn != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, fn)
			}
		})
	}
}

func TestRead(t *testing.T) {
	input := `42
PodQuiz 42
Mon, 5 Aug 1974 00:00:00 +0000
http://mp3.podquiz.com/pq42.mp3
12345
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
	if quiz.URL != "http://mp3.podquiz.com/pq42.mp3" {
		t.Errorf("bad quiz number %s", quiz.URL)
	}
	if quiz.Size != 12345 {
		t.Errorf("bad quiz size %d", quiz.Size)
	}
	if !strings.Contains(quiz.Description, "some lines") {
		t.Errorf("bad descriptions %s", quiz.Description)
	}
}
