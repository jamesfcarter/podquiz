package quiz

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Database represents a flat-file database of quizzes
type Database struct {
	Dir        string
	Cache      map[int]*Episode
	LastUpdate time.Time
	sync.RWMutex
}

// New returns a newly initialized database loaded from dir
func NewDatabase(dir string) *Database {
	d := &Database{
		Dir:   dir,
		Cache: make(map[int]*Episode),
	}
	d.Update()
	return d
}

// Find locates a particular episode by number, returning nil if it does
// not exist.
func (d *Database) Find(episode int) *Episode {
	d.RLock()
	defer d.RUnlock()
	return d.Cache[episode]
}

// MostRecent returns the highest episode number found in the database.
func (d *Database) MostRecent() int {
	d.RLock()
	defer d.RUnlock()
	result := 0
	for i, _ := range d.Cache {
		if i > result {
			result = i
		}
	}
	return result
}

// Count returns the total number of episodes from (and including) from
func (d *Database) Count(from int) int {
	d.RLock()
	defer d.RUnlock()
	result := 0
	for i, _ := range d.Cache {
		if i >= from {
			result += 1
		}
	}
	return result
}

// Page returns an array of quizzes, newest first starting at first and
// containing up to size entries.
func (d *Database) Page(first, size int) []*Episode {
	result := make([]*Episode, 0, size)
	for epno := first; epno > 0; epno-- {
		ep := d.Find(epno)
		if ep == nil {
			continue
		}
		result = append(result, ep)
		if len(result) == size {
			break
		}
	}
	return result
}

// All returns an array of all quizzes, newest first
func (d *Database) All() []*Episode {
	mr := d.MostRecent()
	return d.Page(mr, mr)
}

// Update reconciles any changes made on the on-disc database with the cached
// copy.
func (d *Database) Update() error {
	d.Lock()
	defer d.Unlock()
	updateTime := time.Now()
	for i, q := range d.Cache {
		_, err := os.Stat(q.Filename(d.Dir))
		if err != nil {
			delete(d.Cache, i)
		}
	}
	for _, q := range d.Cache {
		_, err := os.Stat(q.CommentsFilename(d.Dir))
		if err != nil {
			q.Comments = []Comment{}
		}
	}
	files, err := ioutil.ReadDir(d.Dir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		if !fi.ModTime().After(d.LastUpdate) {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".podcast") {
			continue
		}
		n, err := quizNumberFromFilename(fi.Name())
		if err != nil {
			log.Printf("bad filename: %s\n", fi.Name())
			continue
		}
		r, err := d.openFile(fi.Name())
		if err != nil {
			log.Printf("failed to open %s: %v\n", fi.Name(), err)
			continue
		}
		defer r.Close()
		q, err := Read(r)
		if err != nil {
			log.Printf("failed to read %s: %v\n", fi.Name(), err)
			continue
		}
		if q.Number != n {
			log.Printf("mismatched quiz number %d in %s\n", q.Number, fi.Name())
			continue
		}
		d.Cache[n] = q
	}
	for _, fi := range files {
		if !strings.HasSuffix(fi.Name(), ".comments") {
			continue
		}
		n, err := quizNumberFromFilename(fi.Name())
		if err != nil {
			log.Printf("bad filename: %s\n", fi.Name())
			continue
		}
		if d.Cache[n] == nil {
			log.Printf("orphaned commments file %s\n", fi.Name())
		}
		r, err := d.openFile(fi.Name())
		if err != nil {
			log.Printf("failed to open %s: %v\n", fi.Name(), err)
			continue
		}
		defer r.Close()
		err = d.Cache[n].ReadComments(r)
		if err != nil {
			log.Printf("failed to read %s: %v\n", fi.Name(), err)
			continue
		}
	}
	d.LastUpdate = updateTime
	return nil
}

func quizNumberFromFilename(n string) (int, error) {
	nPart := strings.SplitN(n, ".", 2)[0]
	return strconv.Atoi(nPart)
}

func (d *Database) openFile(n string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.Dir, n))
}
