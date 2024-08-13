//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml api.yaml

package api

import (
	"net/http"
	"vdb/pkg/datastore"

	validator "vdb/pkg/registry/validator"
)

func NewHandler(ds *datastore.DataStore, opts ...Option) (http.Handler, error) {
	s := &server{
		ds: ds,
		//vr: validatorRegistry,
	}

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
	ds *datastore.DataStore
	vr *validator.ValidatorRegistry
}

var _ StrictServerInterface = &server{}
