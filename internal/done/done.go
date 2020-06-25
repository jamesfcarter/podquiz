package done

import (
	"io/ioutil"
	"strings"
)

type Done struct {
	file   string
	Rounds map[string]string
}

func New(file string) (*Done, error) {
	d := &Done{file: file}
	return d, d.Load()
}

func (d *Done) Load() error {
	data, err := ioutil.ReadFile(d.file)
	if err != nil {
		return err
	}
	rounds := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			rounds[parts[0]] = parts[1]
		}
	}
	d.Rounds = rounds
	return nil
}
