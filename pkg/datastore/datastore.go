package datastore

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	validator "vdb/pkg/validator/base"
)

type DefaultDriverFunc func(typ common.TypeName) (driver.Driver, error)

type TypeHandler struct {
	dri  driver.Driver
	vali validator.Validator
}

type DataStore struct {
	handlers map[common.TypeName]TypeHandler

	defaultDriverFunc DefaultDriverFunc
}

func NewDataStore(opts ...DataStoreOption) (*DataStore, error) {
	o := &dsOptions{}
	for _, opt := range opts {
		opt(o)
	}

	return &DataStore{
		defaultDriverFunc: o.defaultDriverFunc,
		handlers:          make(map[common.TypeName]TypeHandler),
	}, nil
}

func (d *DataStore) RegisterType(name common.TypeName, vali validator.Validator, opts ...HandlerOption) error {
	o := &handlerOptions{}
	for _, opt := range opts {
		opt(o)
	}

	if o.dri == nil {
		if d.defaultDriverFunc != nil {
			dri, err := d.defaultDriverFunc(name)
			if err != nil {
				return err
			}
			o.dri = dri
		} else {
			return fmt.Errorf("no driver provided")
		}
	}

	d.handlers[name] = TypeHandler{
		dri:  o.dri,
		vali: vali,
	}

	return nil
}

func (d *DataStore) Get(ctx context.Context, typ common.TypeName, id common.TypeID) (Revision, error) {
	handler, ok := d.handlers[typ]
	if !ok {
		return Revision{}, ErrUnknownType{Type: typ}
	}

	rev, err := handler.dri.GetLatest(ctx, id)
	if err != nil {
		return Revision{}, ErrIdNotFound{Type: typ, Id: id}
	}

	return convertRevision(typ, &rev)
}

func (d *DataStore) GetRevision(ctx context.Context, typ common.TypeName, id common.TypeID, revisionId common.RevisionID) (Revision, error) {
	handler, ok := d.handlers[typ]
	if !ok {
		return Revision{}, ErrUnknownType{Type: typ}
	}

	revs, err := handler.dri.GetRevisions(ctx, id)
	if err != nil {
		return Revision{}, ErrIdNotFound{Type: typ, Id: id}
	}

	if int(revisionId) >= len(revs) {
		return Revision{}, ErrRevisionNotFound{Type: typ, Id: id, RevisionID: revisionId}
	}

	return convertRevision(typ, &revs[revisionId])
}

func (d *DataStore) GetRevisionList(ctx context.Context, typ common.TypeName, id common.TypeID) ([]Revision, error) {
	handler, ok := d.handlers[typ]
	if !ok {
		return nil, ErrUnknownType{Type: typ}
	}

	revs, err := handler.dri.GetRevisions(ctx, id)
	if err != nil {
		return nil, ErrIdNotFound{Type: typ, Id: id}
	}

	ret := make([]Revision, len(revs))
	for i, rev := range revs {
		ret[i], err = convertRevision(typ, &rev)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (d *DataStore) Set(ctx context.Context, typ common.TypeName, id common.TypeID, value common.Value) (Revision, error) {
	handler, ok := d.handlers[typ]
	if !ok {
		return Revision{}, fmt.Errorf("type not found: %s", typ)
	}

	rev, err := handler.dri.Set(ctx, id, value)
	if err != nil {
		return Revision{}, err
	}

	return convertRevision(typ, &rev)
}
