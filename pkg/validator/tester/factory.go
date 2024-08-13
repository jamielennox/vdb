package tester

import (
	"context"
	"vdb/pkg/common"
	"vdb/pkg/validator/base"
)

type factory struct {
}

func NewTesterFactory() base.Factory {
	return &factory{}
}

func (f *factory) Build(ctx context.Context, value common.CollectionValue) (base.Validator, error) {
	return &testerValidator{}, nil
}
