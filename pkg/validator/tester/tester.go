package tester

import (
	"context"
	"vdb/pkg/common"
)

type testerValidator struct {
}

//func NewTesterValidator(opts ...Option) (base.Validator, error) {
//	o := options{}
//	for _, opt := range opts {
//		opt(&o)
//	}
//
//	return &testerValidator{}, nil
//}

func (t testerValidator) Validate(ctx context.Context, value common.CollectionValue) error {
	return nil
}
