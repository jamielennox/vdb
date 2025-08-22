//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml api.yaml

package http

import (
	"context"
	"net/http"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"vdb/pkg/datastore"
)

var tracer = otel.Tracer("vdb.api.http")

func otelMiddleware(next strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		var span trace.Span
		ctx, span = tracer.Start(
			ctx,
			operationID,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		return next(ctx, w, r, request)
	}
}

func NewHandler(ds *datastore.DataStore, opts ...Option) (http.Handler, error) {
	s := &server{
		ds: ds,
	}

	o := options{}

	for _, opt := range opts {
		opt(&o)
	}

	si := NewStrictHandlerWithOptions(
		s,
		[]StrictMiddlewareFunc{
			otelMiddleware,
		},
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
}

var _ StrictServerInterface = &server{}
