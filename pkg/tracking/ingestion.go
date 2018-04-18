package tracking

import (
	"context"
	"path"
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

	health := new(Health)
	_, err = s.client.do(ctx, req, health)
	if err != nil {
		return nil, err
	}

	return health, nil
}
