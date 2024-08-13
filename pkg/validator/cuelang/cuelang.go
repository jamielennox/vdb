package cuelang

import (
	"context"
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"fmt"
	"vdb/pkg/common"
	validator "vdb/pkg/validator/base"
)

type cuelangValidator struct {
	ctx *cue.Context
	v   cue.Value
}

func (c *cuelangValidator) Validate(ctx context.Context, value common.CollectionValue) error {
	cueValue := c.ctx.Encode(value)
	unified := c.v.Unify(cueValue)

	if err := unified.Validate(); err != nil {
		fmt.Println("Validation failed:", err)
		return err
	}

	fmt.Println("Validate passed")
	return nil
}

func NewCuelangValidator(schema string) (validator.Validator, error) {
	ctx := cuecontext.New()
	path := "#Schema"
	v := ctx.CompileString(schema).LookupPath(cue.ParsePath(path))

	if err := v.Err(); err != nil {
		return nil, err
	}

	return &cuelangValidator{
		ctx: ctx,
		v:   v,
	}, nil
}
