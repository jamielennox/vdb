package opa

import (
	"context"
	"errors"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"vdb/pkg/authz/base"
	"vdb/pkg/common"
)

type opaAuth struct {
	query rego.PreparedEvalQuery
}

func getOpaInput(targetType string, event common.Event) map[string]any {
	return map[string]any{
		"operation": string(event.Operation),
		"target": map[string]any{
			"name":   event.Target.Name,
			"labels": event.Target.Labels,
			"type":   targetType,
		},
		"subject": map[string]any{
			"user":  event.Subject.UserName,
			"roles": event.Subject.Roles,
		},
	}
}

func (o *opaAuth) Collection(ctx context.Context, event common.Event) error {
	results, err := o.query.Eval(ctx, rego.EvalInput(getOpaInput("collection", event)))
	if err != nil {
		var errw ast.Errors

		switch {
		case errors.As(err, &errw):
			errs := make(ErrOpaFailures, len(errw))

			for i, e := range errw {
				errs[i] = ErrOpaFailure{
					Code:     e.Code,
					Row:      e.Location.Row,
					Filename: e.Location.File,
					Message:  e.Message,
				}
			}
			return errs
		default:
			return err
		}
	}

	if !results.Allowed() {
		fmt.Println("fail, %+v", results)
		return fmt.Errorf("authorization failure")
	}

	return nil
}

func (o *opaAuth) Revision(ctx context.Context, config base.DataConfig, info common.UserInfo) error {
	fmt.Println("Authorize revision", config)
	return nil
}

func NewOpaAuthorizer(ctx context.Context, opts ...func(r *rego.Rego)) (base.Authorizer, error) {
	query, err := rego.New(opts...).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	return &opaAuth{
		query: query,
	}, nil
}
