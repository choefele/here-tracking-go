package tracking

import (
	"context"
	"path"
)

type service struct {
	client *Client
	path   string
}

type Health struct {
	Message string `json:"message,omitempty"`
}

func (s *service) Health(ctx context.Context) (*Health, error) {
	path := path.Join(s.path, "health")
	req, err := s.client.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	health := new(Health)
	_, err = s.client.do(ctx, req, health)
	if err != nil {
		return nil, err
	}

	return health, nil
}
