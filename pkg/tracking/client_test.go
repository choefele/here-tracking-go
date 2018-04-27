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

func TestClientWithParameters(t *testing.T) {
	r := func(ctx context.Context, request *request, response *response) error { return nil }
	c, err := newClientWithParameters(nil, nil, r)

	if err != nil {
		t.Error("Expected no error")
	}
	if got, want := c.BaseURL.String(), "https://tracking.api.here.com"; got != want {
		t.Errorf("Client BaseURL is %v, want %v", got, want)
	}
	if got := c.httpClient; got == nil {
		t.Errorf("Client HTTP client is nil")
	}
	if got := c.authorizedRequest; got == nil {
		t.Errorf("Client authorizedRequest is nil")
	}
}

func TestRequest(t *testing.T) {
	client, mux, teardown := setupTestClient()
	defer teardown()

	mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, `"request"`+"\n")

		fmt.Fprint(w, `"response"`)
	})

	var result string
	err := client.request(
		context.Background(),
		&request{
			path:   "/path",
			method: http.MethodPost,
			body:   "request",
		},
		&response{
			body: &result,
		},
	)

	if err != nil {
		t.Error("Expected no error")
	}

	if got, want := result, "response"; got != want {
		t.Errorf("Request returned %v, want %v", got, want)
	}
}

func TestRequest_error(t *testing.T) {
	client, mux, teardown := setupTestClient()
	defer teardown()

	mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	err := client.request(
		context.Background(),
		&request{
			path:   "/path",
			method: http.MethodPost,
		},
		&response{},
	)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestNewRequest(t *testing.T) {
	c, _ := newClientWithParameters(nil, nil, nil)

	in := map[string]interface{}{
		"a": 3711,
		"b": "2138",
	}
	req, _ := c.newRequest("GET", "path", in, nil)
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
	c, _ := newClientWithParameters(nil, nil, nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.newRequest("GET", ".", &T{}, nil)

	if err == nil {
		t.Error("Expected error to be returned")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c, _ := newClientWithParameters(nil, nil, nil)
	_, err := c.newRequest("GET", ":", nil, nil)

	if err == nil {
		t.Error("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setupTestClient()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.newRequest("GET", ".", nil, nil)
	got := new(foo)
	client.do(context.Background(), req, got)

	want := &foo{"a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, teardown := setupTestClient()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	req, _ := client.newRequest("GET", ".", nil, nil)
	resp, _ := client.do(context.Background(), req, nil)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP 400 error, got %d status code", resp.StatusCode)
	}
}
