package tracking

import (
	"context"
)

type request struct {
	path    string
	method  string
	body    interface{}
	headers map[string]string
}

type response struct {
	body interface{}
}

type requester interface {
	request(ctx context.Context, request *request, response *response) error
}

type requesterFunc func(ctx context.Context, request *request, response *response) error

func (f requesterFunc) request(ctx context.Context, request *request, response *response) error {
	return f(ctx, request, response)
}
