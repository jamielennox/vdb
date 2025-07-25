package memory

import (
	"context"
	"fmt"
	"sync"
	"vdb/pkg/common"
	"vdb/pkg/driver/base"
)

type memoryData struct {
	data    []common.CollectionValue
	transId common.TransactionId
}

type memoryStore struct {
	lock  sync.RWMutex
	store map[common.CollectionId]*memoryData
	//store map[common.CollectionId][]common.CollectionValue
}

func valuesToRevisions(id common.CollectionId, values []common.CollectionValue) []base.Revision {
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

func (m *memoryStore) GetLatest(ctx context.Context, id common.CollectionId) (base.Revision, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v := m.store[id]
	if v == nil {
		return base.Revision{}, fmt.Errorf("id not found")
	}

	if len(v.data) == 0 {
		return base.Revision{}, fmt.Errorf("no revisions")
	}

	return base.Revision{
		Meta: base.Meta{
			Id:       id,
			Revision: common.RevisionID(len(v.data) - 1),
			Version:  base.DefaultVersion,
		},
		Value: v.data[len(v.data)-1],
	}, nil
}

func (m *memoryStore) GetRevisions(ctx context.Context, id common.CollectionId) ([]base.Revision, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v := m.store[id]
	if v == nil {
		return nil, fmt.Errorf("id not found")
	}

	return valuesToRevisions(id, v.data), nil
}

func (m *memoryStore) Set(ctx context.Context, transId common.TransactionId, data ...base.CollectionData) (base.Transaction, error) {
	//func (m *memoryStore) Set(ctx context.Context, id common.CollectionId, value common.CollectionValue) (base.Revision, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.store == nil {
		m.store = make(map[common.CollectionId]*memoryData)
	}

	trans := base.Transaction{
		Id: transId,
	}

	for _, d := range data {
		if m.store[d.Id] == nil {
			m.store[d.Id] = &memoryData{}
		}

		m.store[d.Id].transId = transId
		m.store[d.Id].data = append(m.store[d.Id].data, d.Value)

		trans.Revisions = append(trans.Revisions, base.Revision{
			Meta: base.Meta{
				Id:       d.Id,
				Revision: common.RevisionID(len(m.store[d.Id].data) - 1),
				Version:  base.DefaultVersion,
			},
			Labels: nil,
			Value:  d.Value,
		})
	}

	return trans, nil
}

func NewMemoryStore(opts ...Option) (base.Driver, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &memoryStore{
		store: make(map[common.CollectionId]*memoryData),
	}, nil
}
