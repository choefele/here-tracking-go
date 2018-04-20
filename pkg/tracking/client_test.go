package tracking

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
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

func TestNewRequest_authorization(t *testing.T) {
	c := NewClient()

	token := "token"
	c.AccessToken = &token
	r, _ := c.newRequest("GET", ".", nil)

	testHeader(t, r, "Authorization", "Bearer token")
}

func TestDo(t *testing.T) {
	client, mux, teardown := setupTestServer()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.newRequest("GET", ".", nil)
	got := new(foo)
	client.do(context.Background(), req, got)

	want := &foo{"a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.newRequest("GET", ".", nil)
	resp, err := client.do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}
