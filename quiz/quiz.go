package quiz

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Quiz represents a single episode of PodQuiz
type Quiz struct {
	Number      int
	Name        string
	URL         string
	Released    time.Time
	Size        int64
	Description string
	Comments    []Comment
}

// Comment represents a comment on a quiz
type Comment struct {
	Time    time.Time
	Author  string
	Comment string
}

var (
	BadFormatError = errors.New("bad format")
)

// Filename returns the filename of this quiz in the flat-file database in
// given directory
func (q *Quiz) Filename(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%d.podcast", q.Number))
}

// Read returns a Quiz read from the supplied io.Reader
func Read(r io.Reader) (*Quiz, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	fields := strings.SplitN(string(data), "\n", 6)
	if len(fields) != 6 {
		return nil, BadFormatError
	}
	number, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return nil, err
	}
	released, err := mail.ParseDate(fields[2])
	if err != nil {
		return nil, err
	}
	return &Quiz{
		Number:      number,
		Name:        fields[1],
		Released:    released,
		URL:         fields[3],
		Size:        size,
		Description: fields[5],
	}, nil
}
