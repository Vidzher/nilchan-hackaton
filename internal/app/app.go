package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"nilchan-hackaton/internal/auth"
	"nilchan-hackaton/internal/config"
	"nilchan-hackaton/internal/httpapi/validation"
	"nilchan-hackaton/internal/leaderboard"
	"nilchan-hackaton/internal/llm"
	"nilchan-hackaton/internal/parser"
	"nilchan-hackaton/internal/profile"
	"nilchan-hackaton/internal/quiz"
	quizgen "nilchan-hackaton/internal/quiz/gen"
	"nilchan-hackaton/internal/resource"
	resourcerepo "nilchan-hackaton/internal/resource/repository"
	"nilchan-hackaton/internal/shop"
	"nilchan-hackaton/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	cfg               *config.Config
	router            chi.Router
	storage           *storage.Storage
	resourceProcessor *resource.Processor
}

func New(cfg *config.Config) (*App, error) {
	store, err := storage.New(cfg.StoragePath)
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
		&http.Client{Timeout: cfg.Firecrawl.Timeout})
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("initialize Firecrawl client: %w", err)
	}
	openRouterClient, err := llm.NewOpenRouterClient(
		cfg.OpenRouter.APIKey,
		cfg.OpenRouter.ModelName)
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("initialize OpenRouter client: %w", err)
	}
	quizGenerator, err := quizgen.NewGenerator(openRouterClient)
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("initialize quiz generator: %w", err)
	}
	resourceRepo := resourcerepo.New(store)
	recovered, err := resourceRepo.RecoverProcessing(context.Background())
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("recover interrupted resource processing: %w", err)
	}
	if recovered > 0 {
		fmt.Printf("recovered %d interrupted processing resources\n", recovered)
	}

	a.resourceProcessor = resource.NewProcessor(resourceRepo, quizGenerator)
	resourceService := resource.NewService(
		resourceRepo,
		firecrawlClient,
		a.resourceProcessor,
		cfg.Firecrawl.Timeout)
	resourceHandler := resource.NewHandler(resourceService, validate)

	quizRepo := quiz.NewRepository(store)
	quizService := quiz.NewService(quizRepo)
	quizHandler := quiz.NewHandler(quizService, validate)

	profileRepo := profile.NewRepository(store)
	profileService := profile.NewService(profileRepo)
	profileHandler := profile.NewHandler(profileService, validate)

	shopRepo := shop.NewRepository(store)
	shopService := shop.NewService(shopRepo)
	shopHandler := shop.NewHandler(shopService, validate)

	leaderboardRepo := leaderboard.NewRepository(store)
	leaderboardService := leaderboard.NewService(leaderboardRepo)
	leaderboardHandler := leaderboard.NewHandler(leaderboardService)

	a.registerRoutes(authHandler, resourceHandler, quizHandler, profileHandler, shopHandler, leaderboardHandler)

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

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	log.Printf("HTTP server started on %s", listener.Addr())

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.Serve(listener)
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
	if a.resourceProcessor != nil {
		a.resourceProcessor.Close()
	}
	return a.storage.Close()
}

func (a *App) registerRoutes(
	authHandler *auth.Handler,
	resourceHandler *resource.Handler,
	quizHandler *quiz.Handler,
	profileHandler *profile.Handler,
	shopHandler *shop.Handler,
	leaderboardHandler *leaderboard.Handler,
) {
	a.router.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post("/login", authHandler.HandleLogin())
		r.Post("/register", authHandler.HandleRegister())

		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Post("/resources", resourceHandler.HandleCreate())
			r.Get("/resources", resourceHandler.HandleList())
			r.Get("/resources/{resourceID}", resourceHandler.HandleGet())
			r.Get("/resources/{resourceID}/quiz", quizHandler.HandleGet())
			r.Post("/resources/{resourceID}/quiz/complete", quizHandler.HandleComplete())
			r.Get("/profile", profileHandler.HandleGet())
			r.Patch("/profile/cosmetics", profileHandler.HandleUpdateCosmetics())
			r.Get("/shop", shopHandler.HandleList())
			r.Post("/shop/purchase", shopHandler.HandlePurchase())
			r.Get("/leaderboard", leaderboardHandler.HandleList())
		})
	})
}
