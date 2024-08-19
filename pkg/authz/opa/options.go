package opa

type options struct {
	query      string
	module     string
	moduleName string
}

type Option func(*options)

func WithQuery(query string) Option {
	return func(o *options) {
		o.query = query
	}
}

func WithModuleName(moduleName string) Option {
	return func(o *options) {
		o.moduleName = moduleName
	}
}
