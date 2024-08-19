package cuelang

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	validator "vdb/pkg/validator/base"
)

const DefaultCueLangValidatorName common.ValidatorName = "cuelang"

type factory struct {
}

func (f *factory) Build(ctx context.Context, value common.ValidatorData) (validator.Validator, error) {
	switch d := value.(type) {
	case string:
		return NewCuelangValidator(d)
	}

	return nil, fmt.Errorf("invalid cuelang build data")
}

func NewCuelangFactory() validator.Factory {
	return &factory{}
}
