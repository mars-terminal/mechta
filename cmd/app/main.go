package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/phuslu/log"
	"golang.org/x/sync/errgroup"

	"github.com/mars-terminal/mechta/internal/server/http"
	shortenerService "github.com/mars-terminal/mechta/internal/service/shortener"
	"github.com/mars-terminal/mechta/internal/storage/postgres"
	shortenerStorage "github.com/mars-terminal/mechta/internal/storage/postgres/shortener"
)

type options struct {
	PostgresURL string `long:"postgres-url" env:"POSTGRES_URL" default:"postgres://shortener:postgres-password@localhost:5432/shortener?sslmode=disable" required:"true"`

	HTTPAddr string `long:"http-addr" default:"0.0.0.0:8000" env:"HTTP_ADDR"`

	ShortenerBaseURL string `long:"shortener-base-url" default:"https://example.com" ENV:"SHORTENER_BASE_URL"`
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Msg("failed to parse args")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.NewDataBase(ctx, opts.PostgresURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	server, err := http.NewServer(
		shortenerService.NewService(
			opts.ShortenerBaseURL,
			shortenerStorage.NewStorage(db),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize shortener")
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		defer cancel()

		for {
			select {
			case <-gCtx.Done():
				return context.Canceled
			case <-ctx.Done():
				return context.Canceled
			case s := <-c:
				log.Printf("got stop signal: stopping, signal: %v", s.String())
				return context.Canceled
			}
		}
	})

	g.Go(func() error {
		log.Info().Str("address", fmt.Sprintf("http://%s", opts.HTTPAddr)).Msg("starting server")
		return server.Listen(opts.HTTPAddr)
	})

	g.Go(func() error {
		<-gCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		return server.ShutdownWithContext(ctx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal().Err(err).Msg("some goroutine fails with error")
	} else {
		log.Info().Msg("exiting...")
	}
}
