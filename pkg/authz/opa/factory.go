package opa

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"

	authz "vdb/pkg/authz/base"
	"vdb/pkg/common"
)

const DefaultOpaAuthorizerName common.AuthorizerName = "opa"

type factory struct {
	query      string
	moduleName string
}

func (f *factory) Build(ctx context.Context, value common.AuthorizerData) (authz.Authorizer, error) {
	switch d := value.(type) {
	case string:
		return NewOpaAuthorizer(ctx,
			rego.Query(f.query),
			rego.Module(f.moduleName, d),
			rego.Trace(true),
		)
	}

	return nil, fmt.Errorf("invalid cuelang build data")
}

func NewOpaFactory(opts ...Option) authz.Factory {
	o := &options{
		query:      "data.authz.allow",
		moduleName: "vdb.rego",
	}

	for _, opt := range opts {
		opt(o)
	}

	return &factory{
		query:      o.query,
		moduleName: o.moduleName,
	}
}
