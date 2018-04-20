package tracking

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://tracking.api.here.com"
)

type Client struct {
	httpClient  *http.Client
	BaseURL     url.URL
	AccessToken *string

	Ingestion *IngestionService
}

func NewClient() *Client {
	c, _ := newClientWithParameters(nil, defaultBaseURL)
	return c
}

func newClientWithParameters(httpClient *http.Client, baseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{httpClient: httpClient, BaseURL: *url}
	c.Ingestion = &IngestionService{&service{client: c, path: "/v2"}}

	return c, nil
}

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
func (c *Client) request(ctx context.Context, request *request, response *response) error {
	req, err := c.newRequest(request.method, request.path, request.body, request.headers)
	if err != nil {
		return err
	}

	_, err = c.do(ctx, req, response.body)
	if err != nil {
		return err
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

	if c.AccessToken != nil {
		req.Header.Set("Authorization", "Bearer "+*c.AccessToken)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
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

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
