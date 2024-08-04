package base

import (
	"context"
	"vdb/pkg/common"
)

type Driver interface {
	GetLatest(ctx context.Context, id common.TypeID) (Revision, error)
	GetRevisions(ctx context.Context, id common.TypeID) ([]Revision, error)
	Set(ctx context.Context, id common.TypeID, value common.Value) (Revision, error)
}
