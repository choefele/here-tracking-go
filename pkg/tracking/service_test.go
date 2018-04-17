package tracking

import (
	"context"
	"fmt"
	"net/http"
	"path"
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

	health, err := service.Health(context.Background())
	if err != nil {
		t.Errorf("Service.Health returned error: %v", err)
	} else {
		if got, want := health.Message, "healthy"; got != want {
			t.Errorf("Service.Health returned %v, want %v", got, want)
		}
	}
}
