package tracking

import "testing"

func TestIngestion(t *testing.T) {
	c := NewClient()

	if got := c.Ingestion.client; got == nil {
		t.Errorf("Ingestion service client is nil")
	}
	if got, want := c.Ingestion.path, "/v2"; got != want {
		t.Errorf("Ingestion service path is %v, want %v", got, want)
	}
}
