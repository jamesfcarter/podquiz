package quiz

import (
	"sync"
	"time"
)

// QuizDB represents a flat-file database of quizzes
type QuizDB struct {
	Dir        string
	Cache      map[int]*Quiz
	LastUpdate time.Time
	sync.RWMutex
}

// New returns a newly initialized database loaded from dir
func New(dir string) *QuizDB {
	d := &QuizDB{
		Dir:   dir,
		Cache: make(map[int]*Quiz),
	}
	d.Update()
	return d
}

// Update reconciles any changes made on the on-disc database with the cached
// copy.
func (d *QuizDB) Update() {
	//
}
