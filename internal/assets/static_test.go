package assets_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jamesfcarter/podquiz/internal/assets"
)

func TestAddHandlers(t *testing.T) {
	mux := http.NewServeMux()
	assets.AddHandlers(mux)
	server := httptest.NewServer(mux)
	defer server.Close()
	res, err := http.Get(server.URL + "/img/pqlogo-trans192.png")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("bad status %s", res.Status)
	}
	ct := res.Header.Get("Content-type")
	if ct != "image/png" {
		t.Errorf("bad content type %s", ct)
	}
	img, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(img, assets.MustAsset("static/img/pqlogo-trans192.png")) != 0 {
		t.Error("bad image returned")
	}
}
