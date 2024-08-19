package collection

import (
	"log/slog"
	authz "vdb/pkg/authz/base"
	noopAuthorizer "vdb/pkg/authz/noop"
	"vdb/pkg/common"
	validator "vdb/pkg/validator/base"
	noopValidator "vdb/pkg/validator/noop"
)

type options struct {
	labels common.Labels

	logger *slog.Logger
	vali   validator.Validator
	authz  authz.Authorizer
}

type Option func(*options)

func WithLabel(key, value string) Option {
	return func(o *options) {
		if o.labels == nil {
			o.labels = make(common.Labels)
		}

		o.labels[key] = value
	}
}

func WithAuthorizer(authorizer authz.Authorizer) Option {
	return func(o *options) {
		o.authz = authorizer
	}
}

func WithValidator(vali validator.Validator) Option {
	return func(o *options) {
		o.vali = vali
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func getCollectionOptions(opts ...Option) *options {
	o := &options{
		logger: slog.Default(),
		labels: make(common.Labels),
		vali:   noopValidator.NewNoopValidator(),
		authz:  noopAuthorizer.NewNoopAuthorizer(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
