package tracking

import "testing"

func TestNewClient(t *testing.T) {
	c := NewClient()

	if got, want := c.BaseURL.String(), "https://tracking.api.here.com/"; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got := c.httpClient; got == nil {
		t.Errorf("NewClient HTTP client is nil")
	}
	if got := c.Ingestion; got == nil {
		t.Errorf("NewClient ingestion service is nil")
	}
}
