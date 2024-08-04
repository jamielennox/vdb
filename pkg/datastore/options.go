package datastore

import driver "vdb/pkg/driver/base"

type dsOptions struct {
	defaultDriverFunc DefaultDriverFunc
}

type DataStoreOption func(*dsOptions)

func WithDefaultDriverFunc(f DefaultDriverFunc) DataStoreOption {
	return func(o *dsOptions) {
		o.defaultDriverFunc = f
	}
}

type handlerOptions struct {
	dri driver.Driver
}

type HandlerOption func(*handlerOptions)

func WithDriver(d driver.Driver) HandlerOption {
	return func(o *handlerOptions) {
		o.dri = d
	}
}
