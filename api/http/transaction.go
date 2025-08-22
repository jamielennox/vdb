package http

import (
	"context"
	"github.com/google/uuid"
	"vdb/pkg/collection"
	"vdb/pkg/common"
)

type transactionTableRevision struct {
	Collection string
	Id         string
	RevisionId *uint64
}

type transactionTable struct {
	Id        string
	Revisions map[string]map[string]uint64
}

func (s *server) CreateTransaction(ctx context.Context, request CreateTransactionRequestObject) (CreateTransactionResponseObject, error) {
	//c, err := s.ds.Get(ctx, common.CollectionName("abc"))

	b := *request.Body
	result := make(map[string]map[string]Revision)

	// FIXME: this will conflict if you make a collection called transactions. Split to different data store.
	c, err := s.ds.Get(ctx, common.CollectionName("transactions"))
	if err != nil {
		return nil, err
	}

	u, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	transId := u.String()

	tableData := transactionTable{
		Id:        transId,
		Revisions: make(map[string]map[string]uint64),
	}

	for typ, data := range b {
		for id, _ := range data {
			if tableData.Revisions[typ] == nil {
				tableData.Revisions[typ] = make(map[string]uint64)
			}

			tableData.Revisions[typ][id] = 0
		}
	}

	_, err = c.Set(ctx, nil, collection.CollectionData{
		Id:    common.CollectionId(transId),
		Value: tableData,
	})

	if err != nil {
		return nil, err
	}

	for typ, data := range b {
		c, err := s.ds.Get(ctx, common.CollectionName(typ))
		if err != nil {
			return nil, err
		}

		innerData := []collection.CollectionData{}

		for id, val := range data {
			innerData = append(innerData, collection.CollectionData{
				Id:    common.CollectionId(id),
				Value: val,
			})
		}

		t, err := c.Set(ctx, &transId, innerData...)

		if result[typ] == nil {
			result[typ] = make(map[string]Revision)
		}

		for _, x := range t.Revisions {
			result[typ][string(x.Meta.Id)] = Revision{
				Labels: nil,
				Meta: Meta{
					Id:       TypeId(x.Meta.Id),
					Revision: RevisionId(x.Meta.Revision),
					Type:     typ,
					Version:  int64(x.Meta.Version),
				},
				Value: x.Value,
			}

			tableData.Revisions[typ][string(x.Meta.Id)] = uint64(x.Meta.Revision)
		}
	}

	_, err = c.Set(ctx, nil, collection.CollectionData{
		Id:    common.CollectionId(transId),
		Value: tableData,
	})

	if err != nil {
		// TODO: undo the transaction
		return nil, err
	}

	return CreateTransaction200JSONResponse(result), nil
}
