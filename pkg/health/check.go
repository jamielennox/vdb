package health

import "context"

type CheckFunc func(ctx context.Context) error
