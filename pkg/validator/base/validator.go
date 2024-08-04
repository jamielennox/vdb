package base

import "context"

type Validator interface {
	Validate(ctx context.Context, value map[string]any) error
}
