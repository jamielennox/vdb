package cuelang

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	"vdb/pkg/validator/base"
)

const DefaultCueLangValidatorName common.ValidatorName = "cuelang"

type factory struct {
}

func (f *factory) Build(ctx context.Context, value common.ValidatorData) (base.Validator, error) {
	switch d := value.(type) {
	case string:
		return NewCuelangValidator(d)
	}

	return nil, fmt.Errorf("Invalid cuelang build data")
}

func NewCuelangFactory() base.Factory {
	return &factory{}
}
