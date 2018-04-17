package tracking

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	if got, want := c.BaseURL.String(), "https://tracking.api.here.com"; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got := c.httpClient; got == nil {
		t.Errorf("NewClient HTTP client is nil")
	}
	if got := c.Ingestion; got == nil {
		t.Errorf("NewClient ingestion service is nil")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient()

	in := map[string]interface{}{
		"a": 3711,
		"b": "2138",
	}
	req, _ := c.newRequest("GET", "path", in)
	body, _ := ioutil.ReadAll(req.Body)

	if got, want := req.URL.String(), "https://tracking.api.here.com/path"; got != want {
		t.Errorf("NewRequest URL is %v, want %v", got, want)
	}
	if got, want := string(body), `{"a":3711,"b":"2138"}`+"\n"; got != want {
		t.Errorf("NewRequest body is %v, want %v", got, want)
	}
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest content type is %v, want %v", got, want)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient()

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.newRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient()
	_, err := c.newRequest("GET", ":", nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}
