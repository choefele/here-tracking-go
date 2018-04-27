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

type AdminClient struct {
	*client

	Email       string
	Password    string
	AccessToken *string

	User *UserService
}

func NewAdminClient(email string, password string) *AdminClient {
	return newAdminClientWithParameters(nil, email, password)
}

func newAdminClientWithParameters(baseURL *string, email string, password string) *AdminClient {
	adminClient := &AdminClient{
		Email:    email,
		Password: password,
	}

	client, _ := newClientWithParameters(nil, baseURL, authorizedUserRequesterFunc(adminClient))
	adminClient.client = client
	adminClient.User = &UserService{&service{client: client, path: "/users/v2"}}

	return adminClient
}

func authorizedUserRequesterFunc(c *AdminClient) requesterFunc {
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
