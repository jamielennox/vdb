package datastore

import driver "vdb/pkg/driver/base"

type dsOptions struct {
	driverFactory driver.Factory
}

type DataStoreOption func(*dsOptions)

func WithDriverFactory(f driver.Factory) DataStoreOption {
	return func(o *dsOptions) {
		o.driverFactory = f
	}
}

//type handlerOptions struct {
//	dri driver.Driver
//}
//
//type HandlerOption func(*handlerOptions)
//
//func WithDriver(d driver.Driver) HandlerOption {
//	return func(o *handlerOptions) {
//		o.dri = d
//	}
//}
