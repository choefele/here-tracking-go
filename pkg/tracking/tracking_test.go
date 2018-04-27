package tracking

import (
	"context"
	"net/http"
	"testing"
)

func TestDeviceClient(t *testing.T) {
	c := NewDeviceClient("deviceID", "deviceSecret")

	if got, want := c.DeviceID, "deviceID"; got != want {
		t.Errorf("Device ID is %v, want %v", got, want)
	}
	if got, want := c.DeviceSecret, "deviceSecret"; got != want {
		t.Errorf("Device secret is %v, want %v", got, want)
	}
	if got := c.AccessToken; got != nil {
		t.Errorf("Expected access token to be nil")
	}
	if got := c.client; got == nil {
		t.Errorf("Expected client not to be nil")
	}
	if got := c.Ingestion; got == nil {
		t.Errorf("Expected ingestion service not to be nil")
	}
}

func TestDeviceClientAuthorizedRequest(t *testing.T) {
	client, mux, teardown := setupTestDeviceClient()
	defer teardown()

	token := "access-token"
	client.AccessToken = &token

	mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer access-token")
		w.WriteHeader(http.StatusOK)
	})

	client.authorizedRequest(
		context.Background(),
		&request{path: "/path"},
		&response{},
	)
}

func TestUserClient(t *testing.T) {
	c := NewUserClient("email", "password")

	if got, want := c.Email, "email"; got != want {
		t.Errorf("Email is %v, want %v", got, want)
	}
	if got, want := c.Password, "password"; got != want {
		t.Errorf("Password is %v, want %v", got, want)
	}
	if got := c.AccessToken; got != nil {
		t.Errorf("Expected access token to be nil")
	}
	if got := c.client; got == nil {
		t.Errorf("Expected client not to be nil")
	}
	if got := c.User; got == nil {
		t.Errorf("Expected user service not to be nil")
	}
}

func TestUserClientAuthorizedRequest(t *testing.T) {
	client, mux, teardown := setupTestUserClient()
	defer teardown()

	token := "access-token"
	client.AccessToken = &token

	mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "access-token")
		w.WriteHeader(http.StatusOK)
	})

	client.authorizedRequest(
		context.Background(),
		&request{path: "/path"},
		&response{},
	)
}
