package api

import (
	"context"
	"fmt"
	"vdb/pkg/collection"

	"vdb/pkg/common"
	"vdb/pkg/datastore"
)

func (s *server) GetDataById(ctx context.Context, request GetDataByIdRequestObject) (GetDataByIdResponseObject, error) {
	c, err := s.ds.Get(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetDataById404JSONResponse(RenderErrUnknownType(e)), nil
		default:
			return GetDataById500JSONResponse(RenderServerError(e)), nil
		}
	}

	revision, err := c.Get(ctx, common.CollectionId(request.Id))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetDataById404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return GetDataById404JSONResponse(RenderErrIdNotFound(e)), nil
		default:
			return GetDataById500JSONResponse(RenderServerError(e)), nil
		}
	}

	return GetDataById200JSONResponse(renderRevision(revision)), nil
}

func (s *server) SetData(ctx context.Context, request SetDataRequestObject) (SetDataResponseObject, error) {
	c, err := s.ds.Get(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return SetData404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return SetData404JSONResponse(RenderErrIdNotFound(e)), nil
		default:
			return SetData500JSONResponse(RenderServerError(e)), nil
		}
	}

	transaction, err := c.Set(ctx, nil, collection.CollectionData{
		Id:    common.CollectionId(request.Id),
		Value: request.Body,
	})
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return SetData404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return SetData404JSONResponse(RenderErrIdNotFound(e)), nil
		default:
			return SetData500JSONResponse(RenderServerError(e)), nil
		}
	}

	if len(transaction.Revisions) != 1 {
		return SetData500JSONResponse(RenderServerError(fmt.Errorf("unexpected response count: %d", len(transaction.Revisions)))), nil
	}

	return SetData200JSONResponse(renderRevision(transaction.Revisions[0])), nil
}

func (s *server) GetDataRevisionById(ctx context.Context, request GetDataRevisionByIdRequestObject) (GetDataRevisionByIdResponseObject, error) {
	c, err := s.ds.Get(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetDataRevisionById404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return GetDataRevisionById404JSONResponse(RenderErrIdNotFound(e)), nil
		case collection.ErrRevisionNotFound:
			return GetDataRevisionById404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return GetDataRevisionById500JSONResponse(RenderServerError(e)), nil
		}
	}

	revision, err := c.GetRevision(ctx, common.CollectionId(request.TypId), common.RevisionID(request.RevId))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetDataRevisionById404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return GetDataRevisionById404JSONResponse(RenderErrIdNotFound(e)), nil
		case collection.ErrRevisionNotFound:
			return GetDataRevisionById404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return GetDataRevisionById500JSONResponse(RenderServerError(e)), nil
		}
	}

	return GetDataRevisionById200JSONResponse(renderRevision(revision)), nil
}

func (s *server) ListDataRevisions(ctx context.Context, request ListDataRevisionsRequestObject) (ListDataRevisionsResponseObject, error) {
	c, err := s.ds.Get(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return ListDataRevisions404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return ListDataRevisions404JSONResponse(RenderErrIdNotFound(e)), nil
		case collection.ErrRevisionNotFound:
			return ListDataRevisions404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return ListDataRevisions500JSONResponse(RenderServerError(e)), nil
		}
	}

	// FIXME: This is a single revision
	revs, err := c.GetRevisions(ctx, common.CollectionId(request.TypId))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return ListDataRevisions404JSONResponse(RenderErrUnknownType(e)), nil
		case collection.ErrIdNotFound:
			return ListDataRevisions404JSONResponse(RenderErrIdNotFound(e)), nil
		case collection.ErrRevisionNotFound:
			return ListDataRevisions404JSONResponse(RenderRevisionIdNotFound(e)), nil
		default:
			return ListDataRevisions500JSONResponse(RenderServerError(e)), nil
		}
	}

	ret := make([]Revision, len(revs))
	for i, rev := range revs {
		ret[i] = renderRevision(rev)
	}

	return ListDataRevisions200JSONResponse(ret), nil
}
