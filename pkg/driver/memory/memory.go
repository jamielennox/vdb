package memory

import (
	"context"
	"fmt"
	"sync"
	"vdb/pkg/common"
	"vdb/pkg/driver/base"
)

type memoryStore struct {
	lock  sync.RWMutex
	store map[common.TypeID][]common.Value
}

func valuesToRevisions(id common.TypeID, values []common.Value) []base.Revision {
	ret := make([]base.Revision, len(values))

	for i, val := range values {
		ret[i] = base.Revision{
			Meta: base.Meta{
				Id:       id,
				Revision: common.RevisionID(i),
				Version:  base.DefaultVersion,
			},
			Value: val,
		}
	}

	return ret
}

func (m *memoryStore) GetLatest(ctx context.Context, id common.TypeID) (base.Revision, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v := m.store[id]
	if v == nil {
		return base.Revision{}, fmt.Errorf("id not found")
	}

	if len(v) == 0 {
		return base.Revision{}, fmt.Errorf("no revisions")
	}

	return base.Revision{
		Meta: base.Meta{
			Id:       id,
			Revision: common.RevisionID(len(v) - 1),
			Version:  base.DefaultVersion,
		},
		Value: v[len(v)-1],
	}, nil
}

func (m *memoryStore) GetRevisions(ctx context.Context, id common.TypeID) ([]base.Revision, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v := m.store[id]
	if v == nil {
		return nil, fmt.Errorf("id not found")
	}

	return valuesToRevisions(id, v), nil
}

func (m *memoryStore) Set(ctx context.Context, id common.TypeID, value common.Value) (base.Revision, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.store == nil {
		m.store = make(map[common.TypeID][]common.Value)
	}

	m.store[id] = append(m.store[id], value)

	return base.Revision{
		Meta: base.Meta{
			Id:       id,
			Revision: common.RevisionID(len(m.store[id]) - 1),
			Version:  base.DefaultVersion,
		},
		Value: value,
	}, nil
}

func NewMemoryStore(opts ...Option) (base.Driver, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &memoryStore{
		store: make(map[common.TypeID][]common.Value),
	}, nil
}
