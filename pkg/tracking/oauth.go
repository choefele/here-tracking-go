package tracking

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
)

func parameterString(deviceID string, nonce string, timestamp int64) string {
	var parameterString string

	parameterString += fmt.Sprintf("oauth_consumer_key=%v", deviceID)
	parameterString += fmt.Sprintf("&oauth_nonce=%v", nonce)
	parameterString += "&oauth_signature_method=HMAC-SHA256"
	parameterString += fmt.Sprintf("&oauth_timestamp=%v", timestamp)
	parameterString += "&oauth_version=1.0"
	parameterString += "&realm=IoT"

	return parameterString
}

func baseString(baseURL url.URL, servicePath, parameterString string) string {
	var baseString string

	u, _ := url.Parse(path.Join(servicePath, "/token"))
	baseURLAsString := baseURL.ResolveReference(u).String()

	baseString += "POST"
	baseString += fmt.Sprintf("&%v", url.QueryEscape(baseURLAsString))
	baseString += fmt.Sprintf("&%v", url.QueryEscape(parameterString))

	return baseString
}

func baseSignature(baseString string, deviceSecret string) string {
	mac := hmac.New(sha256.New, []byte(deviceSecret))
	mac.Write([]byte(baseString))
	rawSignature := mac.Sum(nil)
	baseSignature := base64.StdEncoding.EncodeToString(rawSignature)

	return baseSignature
}

func authorizationValue(deviceID string, nonce string, timestamp int64, baseSignature string) string {
	var authorizationValue string

	authorizationValue += "OAuth"
	authorizationValue += " realm=\"IoT\""
	authorizationValue += fmt.Sprintf(", oauth_consumer_key=\"%v\"", deviceID)
	authorizationValue += ", oauth_signature_method=\"HMAC-SHA256\""
	authorizationValue += fmt.Sprintf(", oauth_timestamp=\"%v\"", timestamp)
	authorizationValue += fmt.Sprintf(", oauth_nonce=\"%v\"", nonce)
	authorizationValue += ", oauth_version=\"1.0\""
	authorizationValue += fmt.Sprintf(", oauth_signature=\"%v\"", baseSignature)

	return authorizationValue
}
