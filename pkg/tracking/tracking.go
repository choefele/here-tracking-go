package tracking

import (
	"context"
)

type DeviceClient struct {
	*client

	DeviceID     string
	DeviceSecret string
	AccessToken  *string

	Ingestion *IngestionService
}

func NewDeviceClient(deviceID string, deviceSecret string) *DeviceClient {
	return newDeviceClientWithParameters(nil, deviceID, deviceSecret)
}

func newDeviceClientWithParameters(baseURL *string, deviceID string, deviceSecret string) *DeviceClient {
	deviceClient := &DeviceClient{
		DeviceID:     deviceID,
		DeviceSecret: deviceSecret,
	}

	client, _ := newClientWithParameters(nil, baseURL, authorizedDeviceRequesterFunc(deviceClient))
	deviceClient.client = client
	deviceClient.Ingestion = &IngestionService{&service{client: client, path: "/v2"}}

	return deviceClient
}

func authorizedDeviceRequesterFunc(c *DeviceClient) requesterFunc {
	return requesterFunc(func(ctx context.Context, request *request, response *response) error {
		if c.AccessToken == nil {
			token, err := c.Ingestion.Token(ctx, c.DeviceID, c.DeviceSecret)
			if err != nil {
				return err
			}

			c.AccessToken = &token.AccessToken
		}

		if request.headers == nil {
			request.headers = map[string]string{}
		}
		request.headers["Authorization"] = "Bearer " + *c.AccessToken

		return c.request(ctx, request, response)
	})
}

type UserClient struct {
	*client

	Email       string
	Password    string
	AccessToken *string

	User *UserService
}

func NewUserClient(email string, password string) *UserClient {
	return newUserClientWithParameters(nil, email, password)
}

func newUserClientWithParameters(baseURL *string, email string, password string) *UserClient {
	userClient := &UserClient{
		Email:    email,
		Password: password,
	}

	client, _ := newClientWithParameters(nil, baseURL, authorizedUserRequesterFunc(userClient))
	userClient.client = client
	userClient.User = &UserService{&service{client: client, path: "/users/v2"}}

	return userClient
}

func authorizedUserRequesterFunc(c *UserClient) requesterFunc {
	return requesterFunc(func(ctx context.Context, request *request, response *response) error {
		if c.AccessToken == nil {
			token, err := c.User.Login(ctx, c.Email, c.Password)
			if err != nil {
				return err
			}

			c.AccessToken = &token.AccessToken
		}

		if request.headers == nil {
			request.headers = map[string]string{}
		}
		request.headers["Authorization"] = *c.AccessToken

		return c.request(ctx, request, response)
	})
}
