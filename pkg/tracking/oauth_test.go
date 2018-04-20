package tracking

import (
	"net/url"
	"testing"
)

func TestParameterString(t *testing.T) {
	got := parameterString(
		"9d9c31be-dd5d-40b1-95af-7d5375c39561",
		"0123456789",
		1513634609,
	)
	want := "oauth_consumer_key=9d9c31be-dd5d-40b1-95af-7d5375c39561&oauth_nonce=0123456789&oauth_signature_method=HMAC-SHA256&oauth_timestamp=1513634609&oauth_version=1.0"
	if got != want {
		t.Errorf("Parameter string got %v, want %v", got, want)
	}
}

func TestBaseString(t *testing.T) {
	baseURL, _ := url.Parse("https://tracking.api.here.com")
	got := baseString(
		*baseURL,
		"/v2",
		"oauth_consumer_key=9d9c31be-dd5d-40b1-95af-7d5375c39561&oauth_nonce=0123456789&oauth_signature_method=HMAC-SHA256&oauth_timestamp=1513634609&oauth_version=1.0",
	)
	want := "POST&https%3A%2F%2Ftracking.api.here.com%2Fv2%2Ftoken&oauth_consumer_key%3D9d9c31be-dd5d-40b1-95af-7d5375c39561%26oauth_nonce%3D0123456789%26oauth_signature_method%3DHMAC-SHA256%26oauth_timestamp%3D1513634609%26oauth_version%3D1.0"
	if got != want {
		t.Errorf("Base string got %v, want %v", got, want)
	}
}

func TestBaseSignature(t *testing.T) {
	got := baseSignature(
		"POST&https%3A%2F%2Ftracking.api.here.com%2Fv2%2Ftoken&oauth_consumer_key%3D9d9c31be-dd5d-40b1-95af-7d5375c39561%26oauth_nonce%3D0123456789%26oauth_signature_method%3DHMAC-SHA256%26oauth_timestamp%3D1513634609%26oauth_version%3D1.0",
		"vHrFUhnxo0hxw2VqR5OXBBnvjeTK0T8etmws8HZ9dvw",
	)
	want := "HUm/KJYtAWTIUEUvumh8t9QNmydBNdIv85PnxzHtU8U="
	if got != want {
		t.Errorf("Base signature got %v, want %v", got, want)
	}
}

func TestAuthorizationValue(t *testing.T) {
	got := authorizationValue(
		"9d9c31be-dd5d-40b1-95af-7d5375c39561",
		"0123456789",
		1513634609,
		"HUm/KJYtAWTIUEUvumh8t9QNmydBNdIv85PnxzHtU8U=",
	)

	want := "OAuth realm=\"IoT\",oauth_consumer_key=\"9d9c31be-dd5d-40b1-95af-7d5375c39561\",oauth_nonce=\"0123456789\",oauth_signature_method=\"HMAC-SHA256\",oauth_timestamp=\"1513634609\",oauth_version=\"1.0\",oauth_signature=\"HUm%2FKJYtAWTIUEUvumh8t9QNmydBNdIv85PnxzHtU8U%3D\""
	if got != want {
		t.Errorf("Authorization value got %v, want %v", got, want)
	}
}
