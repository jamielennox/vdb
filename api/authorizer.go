package api

import (
	"context"
	"vdb/pkg/common"
	"vdb/pkg/datastore"
)

func (s *server) GetAuthorizer(ctx context.Context, request GetAuthorizerRequestObject) (GetAuthorizerResponseObject, error) {
	c, err := s.ds.GetAuthorizer(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetAuthorizer404JSONResponse(RenderErrUnknownType(e)), nil
		default:
			return GetAuthorizer500JSONResponse(RenderServerError(e)), nil
		}
	}

	//labels := make(Labels, len(c.Labels))
	//for k, v := range c.Labels {
	//	labels[k] = v
	//}

	return GetAuthorizer200JSONResponse{
		Meta: Meta{
			Id:       request.Type,
			Revision: RevisionId(c.Meta.Revision),
			Type:     TypeName(c.GetName()),
			Version:  int64(c.Meta.Version),
		},
		Value: c.GetConfig(),
	}, nil
}

func (s *server) SetAuthorizer(ctx context.Context, request SetAuthorizerRequestObject) (SetAuthorizerResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
