package tracking

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"time"
)

type IngestionService struct {
	*service
}

type DataRequest struct {
	Position  *Position `json:"position,omitempty"`
	Timestamp Time      `json:"timestamp,omitempty"`
}

type Position struct {
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Accuracy float64 `json:"accuracy,omitempty"`
}

func (s *IngestionService) Send(ctx context.Context, data []*DataRequest) (*Health, error) {
	result := new(Health)
	err := s.client.authorizedClient().request(
		ctx,
		&request{
			path:   path.Join(s.path, "") + "/", // trailing slash is important
			method: "POST",
			body:   data,
		},
		&response{
			body: result,
		},
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type Token struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExpiresIn   Time   `json:"expiresIn,omitempty"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{time.Unix(0, millis*int64(time.Millisecond))}
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Unix() < 0 {
		return nil, errors.New("Time must be after 1 January 1970 00:00:00 UTC")
	}
	time := fmt.Sprintf("%v", t.UnixNano()/int64(time.Millisecond))
	return []byte(time), nil
}

func (s *IngestionService) Token(ctx context.Context, deviceID string, deviceSecret string) (*Token, error) {
	// Authorization header
	timestamp := time.Now().Unix()
	nonce := "0123456789"
	parameterString := parameterString(deviceID, nonce, timestamp)
	baseString := baseString(s.client.BaseURL, s.path, parameterString)
	baseSignature := baseSignature(baseString, deviceSecret)
	authorizationValue := authorizationValue(deviceID, nonce, timestamp, baseSignature)
	headers := map[string]string{"Authorization": authorizationValue}

	result := new(Token)
	err := s.client.request(
		ctx,
		&request{
			path:    path.Join(s.path, "/token"),
			method:  "POST",
			headers: headers,
		},
		&response{
			body: result,
		},
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}
