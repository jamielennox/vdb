package base

import (
	"context"
	"vdb/pkg/common"
)

type Factory interface {
	Build(ctx context.Context, name common.CollectionName, value common.DriverData) (Driver, error)
}

type CollectionData struct {
	Id    common.CollectionId
	Value common.CollectionValue
}

type Driver interface {
	GetLatest(ctx context.Context, id common.CollectionId) (Revision, error)
	GetRevisions(ctx context.Context, id common.CollectionId) ([]Revision, error)
	Set(ctx context.Context, transId common.TransactionId, data ...CollectionData) (Transaction, error)
}
