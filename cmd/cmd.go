package main

import (
	"context"
	"fmt"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/authz/opa"
	"vdb/pkg/common"
	"vdb/pkg/driver/sql"
	"vdb/pkg/factory"
	validator "vdb/pkg/validator/base"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v3"
	"vdb/api"
	"vdb/pkg/datastore"
	"vdb/pkg/driver/memory"
	"vdb/pkg/health"
	"vdb/pkg/validator/cuelang"
)

const sampleCuelang = `
#Schema: {
	name?: string
	age?:  number
}
`

const oldOPA = `
package example.authz

import rego.v1

default allow := false

allow if {
    input.method == "GET"
    input.path == ["salary", input.subject.user]
}

allow if is_admin

is_admin if "admin" in input.subject.groups
`

const sampleOpa = `
package authz

default allow = false

allow {
    input.target.type == "collection"
}
`

func main() {
	cmd := &cli.Command{
		Name:  "serve",
		Usage: "serve",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: "0.0.0.0",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8080,
				Validator: func(i int64) error {
					if i < 0 || i > 65535 {
						return cli.Exit("port must be between 0 and 65535", 1)
					}
					return nil
				},
			},
		},
		Action: serve,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(ctx context.Context, cmd *cli.Command) error {
	r := chi.NewRouter()
	logger := slog.Default()

	db, err := gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			Logger: slogGorm.New(
				slogGorm.WithHandler(logger.Handler()),
			),
		},
	)
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to database", slog.Any("err", err))
		return err
	}

	configDriver, err := memory.NewMemoryStore()
	if err != nil {
		return err
	}

	dataFactory, err := sql.NewSqlDriverFactory(db)
	//dataFactory, err := memory.NewMemoryDriverFactory()
	if err != nil {
		return err
	}

	authzStore, err := memory.NewMemoryStore()
	if err != nil {
		return err
	}

	authzFactory := factory.NewFactory[common.AuthorizerName, common.AuthorizerData, authz.Authorizer](authzStore)
	if err := authzFactory.Register(opa.DefaultOpaAuthorizerName, opa.NewOpaFactory()); err != nil {
		return err
	}

	validatorStore, err := memory.NewMemoryStore()
	if err != nil {
		return err
	}

	validatorFactory := factory.NewFactory[common.ValidatorName, common.ValidatorData, validator.Validator](validatorStore)
	if err := validatorFactory.Register(cuelang.DefaultCueLangValidatorName, cuelang.NewCuelangFactory()); err != nil {
		return err
	}

	ds, err := datastore.NewDataStore(
		configDriver,
		dataFactory,
		authzFactory,
		validatorFactory,
		datastore.WithLogger(logger),
	)
	if err != nil {
		return err
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
		return err
	}

	handler, err := api.NewHandler(ds)
	if err != nil {
		return err
	}

	r.Mount("/api", handler)

	h, err := health.NewHealth()
	if err != nil {
		return err
	}
	r.Mount("/health", h)

	//chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	//	fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
	//	return nil
	//})

	listenAddr := fmt.Sprintf(
		"%s:%d",
		cmd.String("addr"),
		cmd.Int("port"),
	)

	s := &http.Server{
		Handler: r,
		Addr:    listenAddr,
	}

	slog.Info("starting server", slog.String("addr", listenAddr))
	return s.ListenAndServe()
}
