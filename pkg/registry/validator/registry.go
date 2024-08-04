package validator

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	"vdb/pkg/validator/base"
)

type Validator struct {
	store driver.Driver
}

func NewValidatorRegistry(d driver.Driver) (*Validator, error) {
	return &Validator{
		store: d,
	}, nil
}

func (v *Validator) Register(ctx context.Context, name common.TypeID, validator base.Validator) error {
	_, err := v.store.Set(ctx, name, validator)
	return err
}

func (v *Validator) Get(ctx context.Context, name common.TypeID) (base.Validator, error) {
	rev, err := v.store.GetLatest(ctx, name)
	if err != nil {
		return nil, err
	}

	vali, ok := rev.Value.(base.Validator)
	if !ok {
		return nil, fmt.Errorf("invalid validator type: %T", rev.Value)
	}

	return vali, nil
}
