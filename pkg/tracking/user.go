package tracking

import (
	"context"
	"net/http"
	"path"
)

type UserService struct {
	*service
}

func (s *UserService) ListDevices(ctx context.Context) error {
	err := s.client.authorizedClient().request(
		ctx,
		&request{
			path:   path.Join(s.path, "devices"),
			method: http.MethodGet,
		},
		&response{},
	)

	if err != nil {
		return err
	}
	return nil
}

type UserToken struct {
	UserID       string `json:"userId,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	ExpiresIn    Time   `json:"expiresIn,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func (s *UserService) Login(ctx context.Context, email string, password string) (*UserToken, error) {
	result := new(UserToken)
	err := s.client.request(
		ctx,
		&request{
			path:   path.Join(s.path, "login"),
			method: http.MethodPost,
			body: map[string]string{
				"email":    email,
				"password": password,
			},
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
