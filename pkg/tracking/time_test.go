package tracking

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshall(t *testing.T) {
	time := Time{Time: time.Unix(0, 1234567*int64(time.Millisecond))}
	b, _ := json.Marshal(time)

	if got, want := string(b), `1234567`; got != want {
		t.Errorf("JSON got %v, want %v", got, want)
	}
}

func TestMarshall_error(t *testing.T) {
	time := Time{Time: time.Unix(0, 0*int64(time.Millisecond))}
	_, err := json.Marshal(time)

	if err != nil {
		t.Errorf("Expected error not to be nil")
	}
}

func TestUnmarshall(t *testing.T) {
	var ti Time
	err := json.Unmarshal([]byte(`1234567`), &ti)

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	expected := Time{Time: time.Unix(0, 1234567*int64(time.Millisecond))}
	if got, want := ti, expected; got != want {
		t.Errorf("Data got %v, want %v", got, want)
	}
}
