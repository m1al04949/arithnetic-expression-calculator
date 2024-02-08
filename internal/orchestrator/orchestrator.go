package orchestrator

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/storage"
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

	// Router Initiziling
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {

		// r.Post("/expressions", expressions.PostExpression(logger, store)) // Add Expression
		// r.Get("/tasks", tasks.GetTasksList(logger, store))          // Get Tasks List
	})

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
		// logger.Info("failed to start server")
		return err
	}

	// log.Error("server stopped")

	return fmt.Errorf("server is stopped")
}
