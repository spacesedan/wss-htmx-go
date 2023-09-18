package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spacesedan/wss-htmx-go/cmd/internal"
	internaldomain "github.com/spacesedan/wss-htmx-go/internal"
	"github.com/spacesedan/wss-htmx-go/internal/handlers"
	"github.com/spacesedan/wss-htmx-go/internal/hub"
	"github.com/spf13/viper"
)

func main() {
	errC, err := run()
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run() (<-chan error, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	viper, err := internal.NewViper(logger)
	if err != nil {
		logger.Error("Reading Config Failed", slog.String("err", err.Error()))
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewViper")
	}



	srv, err := newServer(ServerConfig{
		logger: logger,
		viper: viper,
	})
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "newServer")
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// shutdown logic
	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal recieved")


		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)


		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown complete")
	}()

	go func() {
		logger.Info("Listening and serving", slog.String("address", viper.GetString("address")))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, nil
}

type ServerConfig struct {
	logger *slog.Logger
	viper   *viper.Viper
}

func newServer(conf ServerConfig) (*http.Server, error) {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RedirectSlashes)

	// Handle static files
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Services
	hub := hub.NewHub(conf.logger)

	// Handler registration
	handlers.NewWssHandler(hub, conf.logger).Register(r)
	handlers.NewViewHandler().Register(r)
	handlers.NewRestHandler().Register(r)

	return &http.Server{
		Handler: r,
		Addr:    conf.viper.GetString("address"),
	}, nil
}
