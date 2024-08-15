package main

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"vdb/pkg/driver/sql"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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

	db, err := gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			//SlowThreshold: time.Second, // Slow SQL threshold
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		panic("failed to connect database")
	}

	sd, err := sql.NewSqlDriverFactory(db)
	if err != nil {
		return err
	}

	cd, err := memory.NewMemoryStore()
	if err != nil {
		return err
	}

	//df, err := memory.NewMemoryDriverFactory()
	//if err != nil {
	//	return err
	//}

	ds, err := datastore.NewDataStore(cd, sd)
	if err != nil {
		return err
	}

	cueFactory := cuelang.NewCuelangFactory()
	if err := ds.Register(cuelang.DefaultCueLangValidatorName, cueFactory); err != nil {
		return err
	}

	if _, err := ds.Set(ctx, "test", "test", cuelang.DefaultCueLangValidatorName, sampleCuelang); err != nil {
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

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

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
