package main

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vdb/pkg/datastore"
	"vdb/pkg/health"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v3"
	_ "go.opentelemetry.io/otel"
	httpapi "vdb/api/http"
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
				Validator: func(i int) error {
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

func mux(ds *datastore.DataStore) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(otelhttp.NewMiddleware("server"))

	h, err := health.NewHealth()
	if err != nil {
		return nil, err
	}
	r.Mount("/health", h)

	handler, err := httpapi.NewHandler(ds)
	if err != nil {
		return nil, err
	}
	r.Mount("/api", handler)

	//chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	//	fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
	//	return nil
	//})

	return r, nil
}

func serve(ctx context.Context, cmd *cli.Command) (err error) {
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	logger := newLogger("vdb")

	ds, err := newDataStore(ctx, logger, cmd)
	if err != nil {
		return
	}

	r, err := mux(ds)
	if err != nil {
		return
	}

	listenAddr := fmt.Sprintf(
		"%s:%d",
		cmd.String("addr"),
		cmd.Int("port"),
	)

	server := &http.Server{
		Handler: r,
		Addr:    listenAddr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		logger.InfoContext(ctx, "starting server", slog.String("addr", listenAddr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return
		}
		err = nil
		logger.InfoContext(ctx, "Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownRelease()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return
	}

	logger.DebugContext(ctx, "Graceful shutdown complete.")
	return nil
}
