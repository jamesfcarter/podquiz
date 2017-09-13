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

// Episode represents a single episode of PodQuiz
type Episode struct {
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
func (q *Episode) Filename(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%d.podcast", q.Number))
}

// CommentsFilename returns the filename of this quiz's comments in
// the flat-file database in given directory
func (q *Episode) CommentsFilename(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%d.comments", q.Number))
}

// Read returns a Quiz read from the supplied io.Reader
func Read(r io.Reader) (*Episode, error) {
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
	return &Episode{
		Number:      number,
		Name:        fields[1],
		Released:    released,
		URL:         fields[3],
		Size:        size,
		Description: fields[5],
	}, nil
}

// ReadComments replaces the quizzes comments with those read from the supplied
// io.Reader
func (q *Episode) ReadComments(r io.Reader) error {
	comments := []Comment{}
	comment := &Comment{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	state := 0
	for _, l := range lines {
		switch state {
		case 0: //read date
			if l != "" {
				t, err := mail.ParseDate(l)
				if err != nil {
					return err
				}
				comment.Time = t
				state = 1
			}
		case 1: //read author
			comment.Author = l
			state = 2
		case 2: //read comment line
			if l == "." {
				comments = append(comments, *comment)
				comment = &Comment{}
				state = 0
			} else {
				comment.Comment += l + "\n"
			}
		}
	}
	q.Comments = comments
	return nil
}
