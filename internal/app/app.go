package app

import (
	"log"
	"net/http"
	"nilchan-hackaton/internal/auth"
	"nilchan-hackaton/internal/config"
	"nilchan-hackaton/internal/domain/pparser"
	"nilchan-hackaton/internal/shared/middlewares"
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
	storage, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Fatalf("storage init error: %v", err.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	a := &App{
		cfg:     cfg,
		router:  router,
		storage: storage,
	}

	firecrawlClient, err := pparser.NewFirecrawlClient(
		cfg.Firecrawl.APIKey,
		cfg.Firecrawl.BaseURL,
		&http.Client{Timeout: cfg.Firecrawl.Timeout},
	)
	if err != nil {
		return nil, err
	}

	a.registerRouter(firecrawlClient)

	return a, nil
}

func (a *App) Run() {
	srv := http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("failed to start server")
	}
}

func (a *App) registerRouter(fcClient *pparser.FirecrawlClient) {
	authRepo := auth.NewAuthRepository(a.storage)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewHandler(authService)

	parserService := pparser.NewParserService(fcClient)
	_ = parserService

	a.router.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		r.Post("/login", authHandler.HandleLogin())
		r.Post("/register", authHandler.HandleRegister())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)
		})
	})
}
