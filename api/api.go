//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml api.yaml

package api

import (
	"context"
	"net/http"
)

func NewHandler(opts ...Option) (http.Handler, error) {
	s := &server{}
	o := options{}

	for _, opt := range opts {
		opt(&o)
	}

	si := NewStrictHandlerWithOptions(
		s,
		[]StrictMiddlewareFunc{},
		StrictHTTPServerOptions{},
	)

	h := HandlerWithOptions(
		si,
		ChiServerOptions{
			BaseURL:          o.baseURL,
			ErrorHandlerFunc: o.errorHandlerFunc,
		},
	)

	return h, nil
}

type server struct {
}

func (s server) GetPing(ctx context.Context, request GetPingRequestObject) (GetPingResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

var _ StrictServerInterface = &server{}
