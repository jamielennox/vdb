package datastore

import (
	"context"
	"log/slog"
	audit "vdb/pkg/audit/base"
	"vdb/pkg/audit/noop"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/collection"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	"vdb/pkg/factory"
	validator "vdb/pkg/validator/base"
)

type AuthzFactory = factory.Factory[common.AuthorizerName, common.AuthorizerData, authz.Authorizer]
type ValidatorFactory = factory.Factory[common.ValidatorName, common.ValidatorData, validator.Validator]

type DataStore struct {
	auditor audit.Auditor
	logger  *slog.Logger

	driver        driver.Driver
	driverFactory driver.Factory

	authzFactory     *AuthzFactory
	validatorFactory *ValidatorFactory
}

func NewDataStore(
	dri driver.Driver,
	factory driver.Factory,
	authzFactory *AuthzFactory,
	validatorFactory *ValidatorFactory,
	opts ...DataStoreOption,
) (*DataStore, error) {
	o := &dsOptions{
		auditor: noop.NewNoopAuditor(),
		logger:  slog.Default(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return &DataStore{
		auditor:          o.auditor,
		driver:           dri,
		driverFactory:    factory,
		logger:           o.logger,
		authzFactory:     authzFactory,
		validatorFactory: validatorFactory,
	}, nil
}

func (d *DataStore) Set(
	ctx context.Context,
	name common.CollectionName,
	//dri common.DriverData,
	validatorName common.ValidatorName,
	validatorData common.ValidatorData,
	authzName common.AuthorizerName,
	authzData common.AuthorizerData,
	opts ...CollectionOption,
) (*collection.Collection, error) {
	aOutput, err := d.authzFactory.Set(ctx, name, authzName, authzData)
	if err != nil {
		return nil, err
	}

	vOutput, err := d.validatorFactory.Set(ctx, name, validatorName, validatorData)
	if err != nil {
		return nil, err
	}

	return d.getCollection(ctx, name, vOutput.Object.Object, aOutput.Object.Object, opts)
}

func (d *DataStore) Get(ctx context.Context, name common.CollectionName, opts ...CollectionOption) (*collection.Collection, error) {
	return d.getCollection(ctx, name, nil, nil, opts)
}

func (d *DataStore) GetValidator(ctx context.Context, name common.CollectionName) (factory.Output[common.ValidatorName, common.ValidatorData, validator.Validator], error) {
	return d.validatorFactory.Get(ctx, name)
}

func (d *DataStore) GetAuthorizer(ctx context.Context, name common.CollectionName) (factory.Output[common.AuthorizerName, common.AuthorizerData, authz.Authorizer], error) {
	return d.authzFactory.Get(ctx, name)
}

func (d *DataStore) getCollection(
	ctx context.Context,
	name common.CollectionName,
	validater validator.Validator,
	authorizer authz.Authorizer,
	opts []CollectionOption,
) (*collection.Collection, error) {
	co := getCollectionOptions(opts...)

	if validater == nil {
		vObj, err := d.GetValidator(ctx, name)
		if err != nil {
			return nil, err
		}

		validater = vObj.Object.Object
	}

	if authorizer == nil {
		aObj, err := d.GetAuthorizer(ctx, name)
		if err != nil {
			return nil, err
		}

		authorizer = aObj.Object.Object
	}

	dri, err := d.driverFactory.Build(ctx, name, "")
	if err != nil {
		return nil, err
	}

	c, err := collection.NewCollection(
		name,
		d.auditor,
		dri,
		collection.WithValidator(validater),
		collection.WithAuthorizer(authorizer),
		collection.WithLogger(d.logger),
	)
	if err != nil {
		return nil, err
	}

	event := common.Event{
		Operation: common.OperationRead,
		Target: common.CollectionTarget{
			Name:   c.Name,
			Type:   "collection",
			Labels: c.Labels,
		},
		Subject: common.UserInfo{
			UserName: "jamie",
			Roles:    []string{"admin"},
		},
	}

	if !co.bypassAuth {
		if err := authorizer.Collection(ctx, event); err != nil {
			return nil, err
		}
	}

	return c, nil
}
