package memory

import (
	"context"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type memoryDriverFactory struct {
	data map[common.CollectionName]driver.Driver
}

func (m *memoryDriverFactory) Build(ctx context.Context, name common.CollectionName, value common.DriverData) (driver.Driver, error) {
	d, ok := m.data[name]

	if !ok {
		d, err := NewMemoryStore()
		if err != nil {
			return nil, err
		}
		m.data[name] = d
	}

	return d, nil
}

func NewMemoryDriverFactory() (driver.Factory, error) {
	return &memoryDriverFactory{
		data: make(map[common.CollectionName]driver.Driver),
	}, nil
}
