package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"vdb/api"
	"vdb/pkg/health"
)

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

	handler, err := api.NewHandler()
	if err != nil {
		return err
	}
	r.Mount("/api", handler)

	h, err := health.NewHealth()
	if err != nil {
		return err
	}
	r.Mount("/health", h)

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
