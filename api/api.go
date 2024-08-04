//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml api.yaml

package api

import (
	"context"
	"net/http"
	"vdb/pkg/common"
	"vdb/pkg/datastore"
)

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
}

func (s *server) ListRevisions(ctx context.Context, request ListRevisionsRequestObject) (ListRevisionsResponseObject, error) {
	revs, err := s.ds.GetRevisionList(ctx, common.TypeName(request.Type), common.TypeID(request.TypId))

	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return ListRevisions404JSONResponse(RenderErrUnknownType(e)), nil
		case datastore.ErrIdNotFound:
			return ListRevisions404JSONResponse(RenderErrIdNotFound(e)), nil
		case datastore.ErrRevisionNotFound:
			return ListRevisions404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return ListRevisions500JSONResponse(RenderServerError(e)), nil
		}
	}

	ret := make([]RevisionSummary, len(revs))

	for i, rev := range revs {
		ret[i] = RevisionSummary{
			Meta: Meta{
				Id:       string(rev.Meta.Id),
				Revision: uint64(rev.Meta.Revision),
				Type:     string(rev.Meta.Type),
				Version:  int64(rev.Meta.Version),
			},
		}
	}

	return ListRevisions200JSONResponse(ret), nil
}

func (s *server) GetRevisionById(ctx context.Context, request GetRevisionByIdRequestObject) (GetRevisionByIdResponseObject, error) {
	revision, err := s.ds.GetRevision(ctx, common.TypeName(request.Type), common.TypeID(request.TypId), common.RevisionID(request.RevId))

	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetRevisionById404JSONResponse(RenderErrUnknownType(e)), nil
		case datastore.ErrIdNotFound:
			return GetRevisionById404JSONResponse(RenderErrIdNotFound(e)), nil
		case datastore.ErrRevisionNotFound:
			return GetRevisionById404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return GetRevisionById500JSONResponse(RenderServerError(e)), nil
		}
	}

	return GetRevisionById200JSONResponse(renderRevision(revision)), nil
}

func (s *server) GetDataTypeId(ctx context.Context, request GetDataTypeIdRequestObject) (GetDataTypeIdResponseObject, error) {
	revision, err := s.ds.Get(ctx, common.TypeName(request.Type), common.TypeID(request.Id))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetDataTypeId404JSONResponse(RenderErrUnknownType(e)), nil
		case datastore.ErrIdNotFound:
			return GetDataTypeId404JSONResponse(RenderErrIdNotFound(e)), nil
		default:
			return GetDataTypeId500JSONResponse(RenderServerError(e)), nil
		}
	}

	return GetDataTypeId200JSONResponse(renderRevision(revision)), nil
}

func (s *server) PutDataTypeId(ctx context.Context, request PutDataTypeIdRequestObject) (PutDataTypeIdResponseObject, error) {
	revision, err := s.ds.Set(ctx, common.TypeName(request.Type), common.TypeID(request.Id), request.Body)
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return PutDataTypeId404JSONResponse(RenderErrUnknownType(e)), nil
		case datastore.ErrIdNotFound:
			return PutDataTypeId404JSONResponse(RenderErrIdNotFound(e)), nil
		default:
			return PutDataTypeId500JSONResponse(RenderServerError(e)), nil
		}
	}

	return PutDataTypeId200JSONResponse(renderRevision(revision)), nil
}

var _ StrictServerInterface = &server{}
