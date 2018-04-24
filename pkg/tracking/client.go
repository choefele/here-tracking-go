package tracking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://tracking.api.here.com"
)

type Client struct {
	httpClient   *http.Client
	BaseURL      url.URL
	AccessToken  *string
	DeviceID     string
	DeviceSecret string

	UserAccessToken *string

	Ingestion *IngestionService
	User      *UserService
}

func NewClient(deviceID string, deviceSecret string) *Client {
	c, _ := newClientWithParameters(nil, defaultBaseURL, deviceID, deviceSecret)
	return c
}

func newClientWithParameters(httpClient *http.Client, baseURL string, deviceID string, deviceSecret string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{httpClient: httpClient, BaseURL: *url, DeviceID: deviceID, DeviceSecret: deviceSecret}
	c.Ingestion = &IngestionService{&service{client: c, path: "/v2"}}
	c.User = &UserService{&service{client: c, path: "/users/v2"}}

	return c, nil
}

func (c *Client) authorizedClient() requester {
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
		if c.UserAccessToken != nil {
			request.headers["Authorization"] = *c.UserAccessToken
		} else {
			request.headers["Authorization"] = "Bearer " + *c.AccessToken
		}

		return c.request(ctx, request, response)
	})
}

func (c *Client) request(ctx context.Context, request *request, response *response) error {
	req, err := c.newRequest(request.method, request.path, request.body, request.headers)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, response.body)
	if err != nil {
		return err
	}

	// Treat response codes != 2xx as error
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP status %v", resp.Status)
	}

	return nil
}

func (c *Client) newRequest(method string, path string, body interface{}, headers map[string]string) (*http.Request, error) {
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	url := c.BaseURL.ResolveReference(pathURL)
	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	if body != nil {
		err = json.NewDecoder(resp.Body).Decode(body)
	}
	return resp, err
}
