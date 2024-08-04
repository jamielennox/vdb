package tester

import (
	"context"

	"vdb/pkg/validator/base"
)

type testerValidator struct {
}

func NewTesterValidator(opts ...Option) (base.Validator, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return &testerValidator{}, nil
}

func (t testerValidator) Validate(ctx context.Context, value map[string]any) error {
	return nil
}
