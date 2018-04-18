package tracking

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"
)

func TestService_Health(t *testing.T) {
	client, mux, teardown := setupTestServer()
	service := client.Ingestion.service
	defer teardown()

	mux.HandleFunc(path.Join(service.path, "health"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"message":"healthy"}`)
	})

	got, err := service.Health(context.Background())
	if err != nil {
		t.Errorf("Service.Health returned error: %v", err)
	}
	want := &Health{Message: "healthy"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response body = %v, want %v", got, want)
	}
}
