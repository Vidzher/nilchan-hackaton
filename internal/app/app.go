package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nilchan-hackaton/internal/auth"
	"nilchan-hackaton/internal/config"
	"nilchan-hackaton/internal/httpapi/validation"
	"nilchan-hackaton/internal/parser"
	"nilchan-hackaton/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	cfg     *config.Config
	router  chi.Router
	storage *storage.Storage
}

func New(cfg *config.Config) (*App, error) {
	store, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("initialize storage: %w", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	a := &App{
		cfg:     cfg,
		router:  router,
		storage: store,
	}

	validate, err := validation.New()
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("initialize validator: %w", err)
	}

	authRepo := auth.NewRepository(store)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService, validate)

	firecrawlClient, err := parser.NewFirecrawlClient(
		cfg.Firecrawl.APIKey,
		cfg.Firecrawl.BaseURL,
		&http.Client{Timeout: cfg.Firecrawl.Timeout},
	)
	if err != nil {
		store.Close()
		return nil, err
	}

	parserService := parser.NewService(firecrawlClient)
	_ = parserService

	a.registerRoutes(authHandler)

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	srv := http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTPServer.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shut down HTTP server: %w", err)
		}

		err := <-serverErrors
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	}
}

func (a *App) Close() error {
	return a.storage.Close()
}

func (a *App) registerRoutes(authHandler *auth.Handler) {
	a.router.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		r.Post("/login", authHandler.HandleLogin())
		r.Post("/register", authHandler.HandleRegister())

		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware)
		})
	})
}
