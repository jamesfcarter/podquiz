package quiz

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Episode represents a single episode of PodQuiz
type Episode struct {
	Number         int
	Name           string
	URL            string
	Released       time.Time
	Size           int64
	RestrictedSize int64
	Description    string
	Comments       []Comment
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

func https(s string) string {
	if strings.HasPrefix(s, "http://") {
		s = "https" + strings.TrimPrefix(s, "http")
	}
	return s
}

// HTML formats the comment for display in the browser
func (c Comment) HTML() template.HTML {
	escaped := template.HTMLEscapeString(c.Comment)
	lines := strings.Split(escaped, "\n")
	return template.HTML(strings.Join(lines, "<br>\n"))
}

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

// DescriptionHTML returns the desciption as a template.HTML
func (q *Episode) DescriptionHTML() template.HTML {
	return template.HTML(q.Description)
}

// OldURL returns an http link to the mp3 for the RSS
func (q *Episode) OldURL() string {
	s := q.URL
	if strings.HasPrefix(s, "https://") {
		s = "http" + strings.TrimPrefix(s, "https")
	}
	return s
}

// RestrictedURL returns an http link to the restricted version of the mp3
func (q *Episode) RestrictedURL() string {
	return strings.ReplaceAll(q.OldURL(), "/pq", "/mpq")
}

// SiteURL returns the absolute or relative URL of the quiz
func (q *Episode) SiteURL(abs bool) string {
	base := "http://podquiz.com"
	if !abs {
		base = ""
	}
	return fmt.Sprintf("%s/quiz.html?q=%d", base, q.Number)
}

// GUID returns a permalink for the episode (using the old URL for consistency)
func (q *Episode) GUID() string {
	return fmt.Sprintf("http://www.podquiz.com/quiz.php?q=pq/%d", q.Number)
}

// MP3 returns the name of the mp3 file for the quiz
func (q *Episode) MP3() string {
	return filepath.Base(q.URL)
}

// Length returns the length of the show as MM:SS
func (q *Episode) Length() string {
	minutes := q.Size / 480000
	seconds := (q.Size - (minutes * 480000)) / 8000
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// CommentCount retuns the number of comments on the quiz
func (q *Episode) CommentCount() int {
	return len(q.Comments)
}

// Read returns a Quiz read from the supplied io.Reader
func Read(r io.Reader) (*Episode, error) {
	data, err := io.ReadAll(r)
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
	sizes := strings.Split(fields[4], ",")
	size, err := strconv.ParseInt(sizes[0], 10, 64)
	if err != nil {
		return nil, err
	}
	var restrictedSize int64
	if len(sizes) > 1 {
		restrictedSize, err = strconv.ParseInt(sizes[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}
	released, err := mail.ParseDate(fields[2])
	if err != nil {
		return nil, err
	}
	return &Episode{
		Number:         number,
		Name:           fields[1],
		Released:       released,
		URL:            https(fields[3]),
		Size:           size,
		RestrictedSize: restrictedSize,
		Description:    fields[5],
	}, nil
}

// ReadComments replaces the quizzes comments with those read from the supplied
// io.Reader
func (q *Episode) ReadComments(r io.Reader) error {
	comments := []Comment{}
	comment := &Comment{}
	data, err := io.ReadAll(r)
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

// AddComment adds the given comment to the on-disc database and to
// the in-memory copy
func (q *Episode) AddComment(dir string, c Comment) error {
	f, err := os.OpenFile(q.CommentsFilename(dir), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	str := c.Time.Format("Mon, 2 Jan 2006 15:04:05 -0700") + "\n" + c.Author + "\n" + c.Comment + "\n.\n"
	_, err = f.WriteString(str)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	q.Comments = append(q.Comments, c)
	return nil
}
