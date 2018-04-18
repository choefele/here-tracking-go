package tracking

import (
	"net/url"
	"testing"
)

func TestParameterString(t *testing.T) {
	got := parameterString(
		"device-id-1234",
		"LIIpk4",
		1513634609,
	)
	want := "oauth_consumer_key=device-id-1234&oauth_nonce=LIIpk4&oauth_signature_method=HMAC-SHA256&oauth_timestamp=1513634609&oauth_version=1.0&realm=IoT"
	if got != want {
		t.Errorf("Parameter string got %v, want %v", got, want)
	}
}

func TestBaseString(t *testing.T) {
	baseURL, _ := url.Parse("https://tracking.api.here.com")
	got := baseString(
		*baseURL,
		"/v2",
		"oauth_consumer_key=device-id-1234&oauth_nonce=LIIpk4&oauth_signature_method=HMAC-SHA256&oauth_timestamp=1513634609&oauth_version=1.0&realm=IoT",
	)
	want := "POST&https%3A%2F%2Ftracking.api.here.com%2Fv2%2Ftoken&oauth_consumer_key%3Ddevice-id-1234%26oauth_nonce%3DLIIpk4%26oauth_signature_method%3DHMAC-SHA256%26oauth_timestamp%3D1513634609%26oauth_version%3D1.0%26realm%3DIoT"
	if got != want {
		t.Errorf("Base string got %v, want %v", got, want)
	}
}

func TestBaseSignature(t *testing.T) {
	got := baseSignature(
		"POST&https%3A%2F%2Ftracking.api.here.com%2Fv2%2Ftoken&oauth_consumer_key%3Ddevice-id-1234%26oauth_nonce%3DLIIpk4%26oauth_signature_method%3DHMAC-SHA256%26oauth_timestamp%3D1513634609%26oauth_version%3D1.0%26realm%3DIoT",
		"device-secret",
	)
	want := "+73cWVgwRPa9gqaO6awyaCjNzlWmVMsLAry8mjlbjdQ="
	if got != want {
		t.Errorf("Base signature got %v, want %v", got, want)
	}
}

func TestAuthorizationValue(t *testing.T) {
	got := authorizationValue(
		"device-id-1234",
		"LIIpk4",
		1513634609,
		"+73cWVgwRPa9gqaO6awyaCjNzlWmVMsLAry8mjlbjdQ=",
	)
	want := "OAuth realm=\"IoT\", oauth_consumer_key=\"device-id-1234\", oauth_signature_method=\"HMAC-SHA256\", oauth_timestamp=\"1513634609\", oauth_nonce=\"LIIpk4\", oauth_version=\"1.0\", oauth_signature=\"+73cWVgwRPa9gqaO6awyaCjNzlWmVMsLAry8mjlbjdQ=\""
	if got != want {
		t.Errorf("Authorization value got %v, want %v", got, want)
	}
}
