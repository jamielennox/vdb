package base

import (
	"context"
	"vdb/pkg/common"
)

type Factory interface {
	Build(ctx context.Context, value common.AuthorizerData) (Authorizer, error)
}

type CollectionConfig struct {
	Name   common.CollectionName
	Labels common.Labels
}

type DataConfig struct {
}

type Authorizer interface {
	Collection(ctx context.Context, event common.Event) error
	Revision(ctx context.Context, data DataConfig, info common.UserInfo) error
}
