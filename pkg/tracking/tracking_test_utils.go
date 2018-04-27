package tracking

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDeviceClient() (*DeviceClient, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	deviceClient := newDeviceClientWithParameters(&server.URL, "", "")
	token := "token"
	deviceClient.AccessToken = &token

	return deviceClient, mux, server.Close
}

func setupTestAdminClient() (*AdminClient, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	adminClient := newAdminClientWithParameters(&server.URL, "", "")
	token := "token"
	adminClient.AccessToken = &token

	return adminClient, mux, server.Close
}

// setup sets up a test HTTP server along with a tracking.client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupTestClient() (*client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ := newClientWithParameters(nil, &server.URL, nil)

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method is %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Request body is %s, want %s", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header value for %q is %q, want %q", header, got, want)
	}
}
