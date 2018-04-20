package tracking

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestIngestion(t *testing.T) {
	c := NewClient("", "")

	if got := c.Ingestion.client; got == nil {
		t.Errorf("Ingestion service client is nil")
	}
	if got, want := c.Ingestion.path, "/v2"; got != want {
		t.Errorf("Ingestion service path is %v, want %v", got, want)
	}
}

func TestIngestion_Send(t *testing.T) {
	client, mux, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc(path.Join(client.Ingestion.path, "")+"/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `[{"position":{"lat":52,"lng":13,"accuracy":100},"timestamp":86399}]`+"\n")
		fmt.Fprint(w, `{"message":"healthy"}`)
	})

	dr := &DataRequest{
		Timestamp: Time{Time: time.Unix(0, 86399*int64(time.Millisecond))},
		Position: &Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	got, err := client.Ingestion.Send(context.Background(), []*DataRequest{dr})
	if err != nil {
		t.Errorf("Ingestion.Send returned error: %v", err)
	}
	want := &Health{Message: "healthy"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}

func TestIngestion_Token(t *testing.T) {
	client, mux, teardown := setupTestServer()
	defer teardown()

	mux.HandleFunc(path.Join(client.Ingestion.path, "/token"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, "")

		if got, want := r.Header.Get("Authorization"), "OAuth realm="; !strings.HasPrefix(got, want) {
			t.Errorf("Header value for \"Authorization\" is %q, want string starting with %q", got, want)
		}

		fmt.Fprint(w, `{"accessToken":"accessToken","expiresIn":86399}`)
	})

	got, err := client.Ingestion.Token(context.Background(), "deviceID", "deviceSecret")
	if err != nil {
		t.Errorf("Ingestion.Token returned error: %v", err)
	}
	want := &Token{AccessToken: "accessToken", ExpiresIn: Time{Time: time.Unix(0, 86399*int64(time.Millisecond))}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}
