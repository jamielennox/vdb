package base

import (
	"context"

	"vdb/pkg/common"
)

type Factory interface {
	Build(ctx context.Context, value common.ValidatorData) (Validator, error)
}

type Validator interface {
	Validate(ctx context.Context, value common.CollectionValue) error
}
