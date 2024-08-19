package datastore

import (
	"log/slog"
	audit "vdb/pkg/audit/base"
	driver "vdb/pkg/driver/base"
)

type dsOptions struct {
	auditor       audit.Auditor
	driverFactory driver.Factory
	logger        *slog.Logger
}

type DataStoreOption func(*dsOptions)

func WithDriverFactory(f driver.Factory) DataStoreOption {
	return func(o *dsOptions) {
		o.driverFactory = f
	}
}

func WithAuditor(a audit.Auditor) DataStoreOption {
	return func(o *dsOptions) {
		o.auditor = a
	}
}

func WithLogger(logger *slog.Logger) DataStoreOption {
	return func(o *dsOptions) {
		o.logger = logger
	}
}

type collectionOptions struct {
	bypassAuth bool
}

type CollectionOption func(options *collectionOptions)

func WithAuthBypass(val bool) CollectionOption {
	return func(o *collectionOptions) {
		o.bypassAuth = val
	}
}

func getCollectionOptions(opts ...CollectionOption) *collectionOptions {
	co := &collectionOptions{
		bypassAuth: false,
	}

	for _, o := range opts {
		o(co)
	}

	return co
}
