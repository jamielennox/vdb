package api

import (
	"context"
	"vdb/pkg/common"
	"vdb/pkg/datastore"
)

func (s *server) GetValidatorSummary(ctx context.Context, request GetValidatorSummaryRequestObject) (GetValidatorSummaryResponseObject, error) {
	c, err := s.ds.GetValidator(ctx, common.CollectionName(request.Type))
	if err != nil {
		switch e := err.(type) {
		case datastore.ErrUnknownType:
			return GetValidatorSummary404JSONResponse(RenderErrUnknownType(e)), nil
		default:
			return GetValidatorSummary500JSONResponse(RenderServerError(e)), nil
		}
	}

	//labels := make(Labels, len(c.Labels))
	//for k, v := range c.Labels {
	//	labels[k] = v
	//}

	return GetValidatorSummary200JSONResponse{
		Meta: Meta{
			Id:       request.Type,
			Revision: RevisionId(c.Meta.Revision),
			Type:     TypeName(c.GetName()),
			Version:  int64(c.Meta.Version),
		},
		Value: c.GetConfig(),
	}, nil
}

func (s *server) SetValidator(ctx context.Context, request SetValidatorRequestObject) (SetValidatorResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
