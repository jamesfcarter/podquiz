package quiz_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jamesfcarter/podquiz/quiz"
)

func TestUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "pqtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "42.podcast"), []byte(`42
PodQuiz 42
Mon, 5 Aug 1974 00:00:00 +0000
http://mp3.podquiz.com/pq42.mp3
12345
Here is a description.
It has some lines.
`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "11.podcast"), []byte(`11
PodQuiz 11
Mon, 5 Aug 1974 00:00:00 +0000
http://mp3.podquiz.com/pq11.mp3
12345
Here is a description.
It has some lines.
`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-20 * time.Minute)
	err = os.Chtimes(filepath.Join(dir, "11.podcast"), oldTime, oldTime)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "42.comments"), []byte(`Mon, 5 Aug 1974 00:00:00 +0000
James
This is a comment.
.
`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	db := &quiz.Database{
		Dir: dir,
		Cache: map[int]*quiz.Episode{
			1: &quiz.Episode{
				Number: 1,
			},
		},
		LastUpdate: time.Now().Add(-10 * time.Minute),
	}
	err = db.Update()
	if err != nil {
		t.Fatal(err)
	}
	if db.Cache[42] == nil {
		t.Error("failed to create episode 42")
	}
	if db.Cache[11] != nil {
		t.Error("failed to ignore old file for episode 11")
	}
	if db.Cache[1] != nil {
		t.Error("failed to remove episode 1")
	}
}

func TestFind(t *testing.T) {
	db := &quiz.Database{
		Cache: map[int]*quiz.Episode{
			1: &quiz.Episode{
				Number: 1,
				Name:   "test",
			},
			2: &quiz.Episode{
				Number: 2,
				Name:   "foo",
			},
		},
	}

	if db.Find(42) != nil {
		t.Error("found missing episode")
	}
	e := db.Find(1)
	if e == nil {
		t.Error("failed to find existing episode")
	}
	if e.Name != "test" {
		t.Error("returned wrong episode")
	}
}

func TestMostRecent(t *testing.T) {
	db := &quiz.Database{
		Cache: map[int]*quiz.Episode{
			1: &quiz.Episode{
				Number: 1,
				Name:   "test",
			},
			42: &quiz.Episode{
				Number: 42,
				Name:   "foo",
			},
		},
	}

	mr := db.MostRecent()
	if mr != 42 {
		t.Errorf("expected 42, got %d", mr)
	}
}
