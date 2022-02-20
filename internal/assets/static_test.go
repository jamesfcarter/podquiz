package assets

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddStaticHandlers(t *testing.T) {
	mux := http.NewServeMux()
	AddStaticHandlers(mux)
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
	if length := len(img); length != 10387 {
		t.Errorf("wrong length %d", length)
	}
}
