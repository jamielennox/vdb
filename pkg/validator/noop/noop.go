package noop

import (
	"context"

	"vdb/pkg/common"
	validator "vdb/pkg/validator/base"
)

type factory struct {
}

type noopValidator struct {
}

func (n *noopValidator) Validate(_ context.Context, _ common.CollectionValue) error {
	return nil
}

func (f *factory) Build(_ context.Context, _ common.ValidatorData) (validator.Validator, error) {
	return NewNoopValidator(), nil
}

func NewNoopValidator() validator.Validator {
	return &noopValidator{}
}

func NewNoopValidatorFactory() validator.Factory {
	return &factory{}
}
