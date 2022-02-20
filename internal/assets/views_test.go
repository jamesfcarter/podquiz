package assets

import (
	"io/ioutil"
	"testing"
)

func TestTemplates(t *testing.T) {
	templates, err := Templates()
	if err != nil {
		t.Fatal(err)
	}
	template := templates["guide.html"]
	err = template.Execute(ioutil.Discard, nil)
	if err != nil {
		t.Error(err)
	}
}
