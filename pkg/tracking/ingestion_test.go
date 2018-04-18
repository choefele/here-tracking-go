package tracking

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"
)

func TestIngestion(t *testing.T) {
	c := NewClient()

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
		testBody(t, r, `{"position":{"lat":52,"lng":13,"accuracy":100}}`+"\n")
		fmt.Fprint(w, `{"message":"healthy"}`)
	})

	dr := &DataRequest{
		Position: &Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	got, err := client.Ingestion.Send(context.Background(), dr)
	if err != nil {
		t.Errorf("Ingestion.Send returned error: %v", err)
	}
	want := &Health{Message: "healthy"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}
