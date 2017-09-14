package assets_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func TestExec(t *testing.T) {
	data := assets.LayoutTemplateData{
		PageTitle: "Test Title",
	}
	templates, err := assets.MakeTemplates()
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	err = templates.Exec("guide", buf, data)
	if err != nil {
		t.Fatal(err)
	}
	result := buf.String()
	if !strings.Contains(result, "<title>Test Title</title>") {
		t.Errorf("document did not contain expected Test Title\n%s", result)
	}
}
