package main

import (
	"context"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/urfave/cli/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/authz/opa"
	"vdb/pkg/common"
	"vdb/pkg/datastore"
	"vdb/pkg/driver/memory"
	"vdb/pkg/driver/sql"
	"vdb/pkg/factory"
	validator "vdb/pkg/validator/base"
	"vdb/pkg/validator/cuelang"
)

func newDataStore(ctx context.Context, logger *slog.Logger, cmd *cli.Command) (*datastore.DataStore, error) {
	db, err := gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			Logger: slogGorm.New(
				slogGorm.WithHandler(logger.Handler()),
			),
		},
	)

	if err != nil {
		//logger.ErrorContext(ctx, "failed to connect to database", slog.Any("err", err))
		return nil, err
	}

	configDriver, err := memory.NewMemoryStore()
	if err != nil {
		return nil, err
	}

	dataFactory, err := sql.NewSqlDriverFactory(db)
	//dataFactory, err := memory.NewMemoryDriverFactory()
	if err != nil {
		return nil, err
	}

	authzStore, err := memory.NewMemoryStore()
	if err != nil {
		return nil, err
	}

	authzFactory := factory.NewFactory[common.AuthorizerName, common.AuthorizerData, authz.Authorizer](authzStore)
	if err := authzFactory.Register(opa.DefaultOpaAuthorizerName, opa.NewOpaFactory()); err != nil {
		return nil, err
	}

	validatorStore, err := memory.NewMemoryStore()
	if err != nil {
		return nil, err
	}

	validatorFactory := factory.NewFactory[common.ValidatorName, common.ValidatorData, validator.Validator](validatorStore)
	if err := validatorFactory.Register(cuelang.DefaultCueLangValidatorName, cuelang.NewCuelangFactory()); err != nil {
		return nil, err
	}

	ds, err := datastore.NewDataStore(
		configDriver,
		dataFactory,
		authzFactory,
		validatorFactory,
		datastore.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	if _, err := ds.Set(
		ctx,
		"test",
		cuelang.DefaultCueLangValidatorName,
		sampleCuelang,
		opa.DefaultOpaAuthorizerName,
		sampleOpa,
		datastore.WithAuthBypass(true),
	); err != nil {
		return nil, err
	}

	return ds, nil
}
