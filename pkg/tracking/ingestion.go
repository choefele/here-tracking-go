package tracking

import (
	"context"
	"path"
	"strconv"
	"time"
)

type IngestionService struct {
	*service
}

type DataRequest struct {
	Position *Position `json:"position,omitempty"`
}

type Position struct {
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Accuracy float64 `json:"accuracy,omitempty"`
}

func (s *IngestionService) Send(ctx context.Context, data *DataRequest) (*Health, error) {
	path := path.Join(s.path, "") + "/" // trailing slash is important
	req, err := s.client.newRequest("POST", path, data)
	if err != nil {
		return nil, err
	}

	body := new(Health)
	_, err = s.client.do(ctx, req, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type Token struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExpiresIn   Time   `json:"expiresIn,omitempty"`
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}

func (s *IngestionService) Token(ctx context.Context, deviceID string, deviceSecret string) (*Token, error) {
	path := path.Join(s.path, "/token")
	req, err := s.client.newRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	// Authorization
	timestamp := time.Now().Unix()
	nonce := "0123456789"
	parameterString := parameterString(deviceID, nonce, timestamp)
	baseString := baseString(*s.client.BaseURL, s.path, parameterString)
	baseSignature := baseSignature(baseString, deviceSecret)
	authorizationValue := authorizationValue(deviceID, nonce, timestamp, baseSignature)
	req.Header.Set("Authorization", authorizationValue)

	body := new(Token)
	_, err = s.client.do(ctx, req, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
