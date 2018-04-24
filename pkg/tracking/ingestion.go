package tracking

import (
	"context"
	"net/http"
	"path"
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

func (s *IngestionService) Send(ctx context.Context, data []*DataRequest) error {
	result := new(Health)
	err := s.client.authorizedClient().request(
		ctx,
		&request{
			path:   path.Join(s.path, "") + "/", // trailing slash is important
			method: http.MethodPost,
			body:   data,
		},
		&response{
			body: result,
		},
	)

	if err != nil {
		return err
	}
	return nil
}

type Token struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExpiresIn   Time   `json:"expiresIn,omitempty"`
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
			method:  http.MethodPost,
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
