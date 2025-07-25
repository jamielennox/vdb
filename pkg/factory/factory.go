package factory

import (
	"context"
	"fmt"

	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type Builder[V any, O any] interface {
	Build(ctx context.Context, value V) (O, error)
}

type Factory[K ~string, V any, O any] struct {
	driver    driver.Driver
	instances map[K]Builder[V, O]
}

func NewFactory[K ~string, V any, O any](dri driver.Driver) *Factory[K, V, O] {
	return &Factory[K, V, O]{
		driver:    dri,
		instances: make(map[K]Builder[V, O]),
	}
}

func (f *Factory[K, V, O]) Register(name K, builder Builder[V, O]) error {
	f.instances[name] = builder
	return nil
}

func (f *Factory[K, V, O]) Create(ctx context.Context, name K, config V) (Object[K, V, O], error) {
	obj := Object[K, V, O]{
		config: configData[K, V]{
			Name:   name,
			Config: config,
		},
	}

	b, ok := f.instances[name]
	if !ok {
		return obj, fmt.Errorf("builder not found: %s", name)
	}

	o, err := b.Build(ctx, config)
	if err != nil {
		return obj, fmt.Errorf("failed to build object: %w", err)
	}

	obj.Object = o
	return obj, nil
}

func (f *Factory[K, V, O]) Set(ctx context.Context, collection common.CollectionName, name K, value V) (Output[K, V, O], error) {
	o, err := f.Create(ctx, name, value)
	if err != nil {
		return Output[K, V, O]{}, err
	}

	t, err := f.driver.Set(ctx, nil, driver.CollectionData{
		Id:    common.CollectionId(collection),
		Value: o.config,
	})
	if err != nil {
		return Output[K, V, O]{}, fmt.Errorf("failed to set collection %s: %w", collection, err)
	}

	if len(t.Revisions) != 1 {
		return Output[K, V, O]{}, fmt.Errorf("unexpected collection response %s: found(%d), expected(1)", collection, len(t.Revisions))
	}

	return Output[K, V, O]{
		Object:     o,
		Meta:       t.Revisions[0].Meta,
		driver:     f.driver,
		collection: collection,
	}, nil
}

func (f *Factory[K, V, O]) Get(ctx context.Context, name common.CollectionName) (Output[K, V, O], error) {
	out := Output[K, V, O]{
		Object:     Object[K, V, O]{},
		Meta:       driver.Meta{},
		driver:     f.driver,
		collection: name,
	}
	rev, err := f.driver.GetLatest(ctx, common.CollectionId(name))
	if err != nil {
		return out, err
	}
	out.Meta = rev.Meta

	config, ok := rev.Value.(configData[K, V])
	if !ok {
		return out, fmt.Errorf("collection config not found: %s", name)
	}
	out.Object.config = config

	o, err := f.Create(ctx, config.Name, config.Config)
	if err != nil {
		return out, err
	}

	out.Object = o
	return out, nil
}
