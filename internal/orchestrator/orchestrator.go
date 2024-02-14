package orchestrator

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/http-server/handlers/expressions"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/http-server/handlers/pages"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/orchrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/pagesrepository"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/templates"
)

func RunServer() error {

	// Config Initializing
	cfg := config.LoadCfg()

	// Log Initializing
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("log init")

	// Storage Initializing
	store := storage.New(cfg.StoragePath, cfg.DatabaseURL)
	if err := store.Open(); err != nil {
		logger.Error("failed to init storage")
		return err
	}
	defer store.Close()
	if err := store.CreateTabs(); err != nil {
		logger.Error("failed to init tabs")
		return err
	}
	logger.Info("storage is initialized")

	// Templates Initializing
	templates, err := templates.New(*cfg)
	if err != nil {
		logger.Error("failed to init templates")
		return err
	}
	logger.Info("templates is initialized")

	//Init Repositories
	orchRepository := orchrepository.New(logger, store)
	pagesRepository := pagesrepository.New(logger, templates, cfg, store)

	// Init Handlers
	expHandler := expressions.New(*orchRepository)
	pageHandler := pages.New(*pagesRepository)

	// Router Initiziling
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/", pageHandler.GetMainPage)               // Get Main Page
		r.Post("/", expHandler.PostExpression)            // Add Expression
		r.Get("/settings", pageHandler.GetSettingsPage)   // Get Settings Page
		r.Post("/settings", pageHandler.SetSettingsPage)  // Post Settings Page
		r.Get("/expressions", pageHandler.GetExpressions) // Get All Expressions
		// r.Get("/tasks", pageHandler.GetTasks)             // Get Tasks List
	})

	// Check new expressions, parsing and calculate
	done := make(chan struct{})
	go orchRepository.Processing(logger, cfg.ProcessingInterval, done)

	// Start HTTP Server
	logger.Info("starting server address", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Info("failed to start server")
		return err
	}

	logger.Error("server stopped")

	return fmt.Errorf("server is stopped")
}
