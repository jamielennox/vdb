package noop

import (
	"context"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/common"
)

type authorizer struct{}

func (a authorizer) Collection(_ context.Context, _ common.Event) error {
	return nil
}

func (a authorizer) Revision(_ context.Context, _ authz.DataConfig, _ common.UserInfo) error {
	return nil
}

type factory struct{}

func (f factory) Build(_ context.Context, _ common.AuthorizerData) (authz.Authorizer, error) {
	return NewNoopAuthorizer(), nil
}

func NewNoopAuthorizer() authz.Authorizer {
	return &authorizer{}
}

func NewNoopAuthorizerFactory() authz.Factory {
	return &factory{}
}
