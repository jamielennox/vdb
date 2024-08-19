package datastore

import (
	"context"
	"fmt"
	"log/slog"
	audit "vdb/pkg/audit/base"
	"vdb/pkg/audit/noop"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/collection"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	validator "vdb/pkg/validator/base"
)

type DataStore struct {
	auditor audit.Auditor
	logger  *slog.Logger

	driver        driver.Driver
	driverFactory driver.Factory

	authzFactories     map[common.AuthorizerData]authz.Factory
	validatorFactories map[common.ValidatorName]validator.Factory
}

type collectionConfig struct {
	AuthzName       common.AuthorizerName `json:"authz_name"`
	AuthzConfig     common.AuthorizerData `json:"authz_config"`
	ValidatorName   common.ValidatorName  `json:"validator_name"`
	ValidatorConfig common.ValidatorData  `json:"validator_config"`
	DriverConfig    common.DriverData     `json:"driver_config"`
}

func NewDataStore(driver driver.Driver, factory driver.Factory, opts ...DataStoreOption) (*DataStore, error) {
	o := &dsOptions{
		auditor: noop.NewNoopAuditor(),
		logger:  slog.Default(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return &DataStore{
		auditor:            o.auditor,
		driver:             driver,
		driverFactory:      factory,
		logger:             o.logger,
		authzFactories:     make(map[common.AuthorizerData]authz.Factory),
		validatorFactories: make(map[common.ValidatorName]validator.Factory),
	}, nil
}

func (d *DataStore) RegisterAuthorizer(name common.AuthorizerName, factory authz.Factory) error {
	d.authzFactories[name] = factory
	return nil
}

func (d *DataStore) RegisterValidator(name common.ValidatorName, factory validator.Factory) error {
	d.validatorFactories[name] = factory
	return nil
}

func (d *DataStore) Set(
	ctx context.Context,
	name common.CollectionName,
	dri common.DriverData,
	validatorName common.ValidatorName,
	validatorData common.ValidatorData,
	authzName common.AuthorizerName,
	authzData common.AuthorizerData,
	opts ...CollectionOption,
) (*collection.Collection, error) {
	co := getCollectionOptions(opts...)

	config := collectionConfig{
		ValidatorName:   validatorName,
		ValidatorConfig: validatorData,
		AuthzName:       authzName,
		AuthzConfig:     authzData,
		DriverConfig:    dri,
	}

	collection, err := d.getCollection(ctx, name, &config, co)
	if err != nil {
		return nil, err
	}

	if _, err := d.driver.Set(ctx, common.CollectionId(name), config); err != nil {
		return nil, err
	}

	return collection, nil
}

func (d *DataStore) Get(ctx context.Context, name common.CollectionName, opts ...CollectionOption) (*collection.Collection, error) {
	co := getCollectionOptions(opts...)

	rev, err := d.driver.GetLatest(ctx, common.CollectionId(name))
	if err != nil {
		return nil, err
	}

	config, ok := rev.Value.(collectionConfig)
	if !ok {
		return nil, fmt.Errorf("collection config not found: %s", name)
	}

	return d.getCollection(ctx, name, &config, co)
}

func (d *DataStore) getCollection(ctx context.Context, name common.CollectionName, config *collectionConfig, opts *collectionOptions) (*collection.Collection, error) {
	azFactory, ok := d.authzFactories[config.AuthzName]
	if !ok {
		return nil, fmt.Errorf("authorizer factory not found: %s", config.AuthzName)
	}

	validatorFactory, ok := d.validatorFactories[config.ValidatorName]
	if !ok {
		return nil, fmt.Errorf("validator factory not found: %s", config.ValidatorName)
	}

	az, err := azFactory.Build(ctx, config.AuthzConfig)
	if err != nil {
		return nil, err
	}

	v, err := validatorFactory.Build(ctx, config.ValidatorConfig)
	if err != nil {
		return nil, err
	}

	dri, err := d.driverFactory.Build(ctx, name, config.DriverConfig)
	if err != nil {
		return nil, err
	}

	c, err := collection.NewCollection(
		name,
		d.auditor,
		dri,
		collection.WithValidator(v),
		collection.WithAuthorizer(az),
		collection.WithLogger(d.logger),
	)
	if err != nil {
		return nil, err
	}

	event := common.Event{
		Operation: common.OperationRead,
		Target: common.CollectionTarget{
			Name:   c.Name,
			Labels: c.Labels,
		},
		Subject: common.UserInfo{
			UserName: "jamie",
			Roles:    []string{"admin"},
		},
	}

	if !opts.bypassAuth {
		if err := az.Collection(ctx, event); err != nil {
			return nil, err
		}
	}

	return c, nil
}
