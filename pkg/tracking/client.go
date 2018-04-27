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

type client struct {
	httpClient        *http.Client
	authorizedRequest requesterFunc
	BaseURL           url.URL
}

func newClientWithParameters(httpClient *http.Client, baseURL *string, authorizedRequest requesterFunc) (*client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	var baseURLForParsing = defaultBaseURL
	if baseURL != nil {
		baseURLForParsing = *baseURL
	}

	url, err := url.Parse(baseURLForParsing)
	if err != nil {
		return nil, err
	}

	c := &client{httpClient: httpClient, BaseURL: *url, authorizedRequest: authorizedRequest}
	return c, nil
}

func (c *client) request(ctx context.Context, request *request, response *response) error {
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

func (c *client) newRequest(method string, path string, body interface{}, headers map[string]string) (*http.Request, error) {
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

func (c *client) do(ctx context.Context, req *http.Request, body interface{}) (*http.Response, error) {
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
