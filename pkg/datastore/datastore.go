package datastore

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	validator "vdb/pkg/validator/base"
)

type DataStore struct {
	driver             driver.Driver
	driverFactory      driver.Factory
	validatorFactories map[common.ValidatorName]validator.Factory
}

type collectionConfig struct {
	ValidatorName   common.ValidatorName `json:"validator_name"`
	ValidatorConfig common.ValidatorData `json:"validator_config"`
	DriverConfig    common.DriverData    `json:"driver_config"`
}

func NewDataStore(driver driver.Driver, factory driver.Factory, opts ...DataStoreOption) (*DataStore, error) {
	o := &dsOptions{}
	for _, opt := range opts {
		opt(o)
	}

	return &DataStore{
		driver:             driver,
		driverFactory:      factory,
		validatorFactories: make(map[common.ValidatorName]validator.Factory),
	}, nil
}

func (d *DataStore) Register(name common.ValidatorName, factory validator.Factory) error {
	d.validatorFactories[name] = factory
	return nil
}

func (d *DataStore) Set(ctx context.Context, name common.CollectionName, dri common.DriverData, vali common.ValidatorName, data common.ValidatorData) (*Collection, error) {
	factory, ok := d.validatorFactories[vali]
	if !ok {
		return nil, fmt.Errorf("Validator factory not found: %s", vali)
	}

	validator, err := factory.Build(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("CollectionValue does not load for factor %s: %w", vali, err)
	}

	driver, err := d.driverFactory.Build(ctx, name, dri)
	if err != nil {
		return nil, fmt.Errorf("Can't load driver: %w", err)
	}

	c := collectionConfig{
		ValidatorName:   vali,
		ValidatorConfig: data,
		DriverConfig:    dri,
	}

	if _, err := d.driver.Set(ctx, common.CollectionId(name), c); err != nil {
		return nil, err
	}

	return &Collection{
		Name: name,
		dri:  driver,
		vali: validator,
	}, nil
}

func (d *DataStore) Get(ctx context.Context, name common.CollectionName) (*Collection, error) {
	rev, err := d.driver.GetLatest(ctx, common.CollectionId(name))
	if err != nil {
		return nil, err
	}

	config, ok := rev.Value.(collectionConfig)
	if !ok {
		return nil, fmt.Errorf("Collection config not found: %s", name)
	}

	vali, ok := d.validatorFactories[config.ValidatorName]
	if !ok {
		return nil, fmt.Errorf("Validator factory not found: %s", config.ValidatorName)
	}

	v, err := vali.Build(ctx, config.ValidatorConfig)
	if err != nil {
		return nil, err
	}

	dri, err := d.driverFactory.Build(ctx, name, config.DriverConfig)
	if err != nil {
		return nil, err
	}

	return &Collection{
		Name: name,
		vali: v,
		dri:  dri,
	}, nil
}

//func (d *DataStore) Register(name common.CollectionName, config CollectionConfig) (*Collection, error) {
//	vali, err := d.registry.Get(config.ValidatorName)
//	if err != nil {
//		return err
//	}
//
//	dri, err := d.defaultDriverFunc(name)
//	if err != nil {
//		return err
//	}
//
//	d.collections[name] = &Collection{
//		typ:  name,
//		dri:  dri,
//		vali: vali,
//	}
//
//	return nil
//}
//
//func (d *DataStore) RegisterType(name common.CollectionName, vali common.ValidatorName, opts ...HandlerOption) error {
//	o := &handlerOptions{}
//	for _, opt := range opts {
//		opt(o)
//	}
//
//	if o.dri == nil {
//		if d.defaultDriverFunc != nil {
//			dri, err := d.defaultDriverFunc(name)
//			if err != nil {
//				return err
//			}
//			o.dri = dri
//		} else {
//			return fmt.Errorf("no driver provided")
//		}
//	}
//
//	d.handlers[name] = TypeHandler{
//		dri:  o.dri,
//		vali: vali,
//	}
//
//	return nil
//}

//func (d *DataStore) Get(ctx context.Context, typ common.CollectionName, id common.CollectionId) (Revision, error) {
//	handler, ok := d.handlers[typ]
//	if !ok {
//		return Revision{}, ErrUnknownType{Type: typ}
//	}
//
//	rev, err := handler.dri.GetLatest(ctx, id)
//	if err != nil {
//		return Revision{}, ErrIdNotFound{Type: typ, Id: id}
//	}
//
//	return convertRevision(typ, &rev)
//}
//
//func (d *DataStore) GetRevision(ctx context.Context, typ common.CollectionName, id common.CollectionId, revisionId common.RevisionID) (Revision, error) {
//	handler, ok := d.handlers[typ]
//	if !ok {
//		return Revision{}, ErrUnknownType{Type: typ}
//	}
//
//	revs, err := handler.dri.GetRevisions(ctx, id)
//	if err != nil {
//		return Revision{}, ErrIdNotFound{Type: typ, Id: id}
//	}
//
//	if int(revisionId) >= len(revs) {
//		return Revision{}, ErrRevisionNotFound{Type: typ, Id: id, RevisionID: revisionId}
//	}
//
//	return convertRevision(typ, &revs[revisionId])
//}
//
//func (d *DataStore) GetRevisionList(ctx context.Context, typ common.CollectionName, id common.CollectionId) ([]Revision, error) {
//	handler, ok := d.handlers[typ]
//	if !ok {
//		return nil, ErrUnknownType{Type: typ}
//	}
//
//	revs, err := handler.dri.GetRevisions(ctx, id)
//	if err != nil {
//		return nil, ErrIdNotFound{Type: typ, Id: id}
//	}
//
//	ret := make([]Revision, len(revs))
//	for i, rev := range revs {
//		ret[i], err = convertRevision(typ, &rev)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return ret, nil
//}
