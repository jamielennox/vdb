package factory

type configData[K ~string, V any] struct {
	Name   K `json:"name"`
	Config V `json:"config"`
}

type Object[K ~string, V, O any] struct {
	config configData[K, V]
	Object O
}

func (o *Object[K, V, O]) GetName() K {
	return o.config.Name
}

func (o *Object[K, V, O]) GetConfig() V {
	return o.config.Config
}
