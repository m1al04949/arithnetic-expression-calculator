package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m1al04949/arithnetic-expression-calculator/internal/config"
	"github.com/m1al04949/arithnetic-expression-calculator/internal/orchestrator"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// Config initialize
	cfg := config.MustLoadCfg()

	// Logger initialize
	log := setupLogger(cfg.Env)

	// Init Orchestrator
	orch := orchestrator.New(cfg, log)

	// Creating channel for graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Run Server
	go func() {
		if err := orch.RunServer(); err != nil {
			log.Error("error server running", slog.Any("error", err))
		}
		log.Info("starting application",
			slog.Any("cfg", cfg))
	}()

	// Given signal for shutdown
	sig := <-sigint
	log.Info("Received", slog.Any("signal", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := orch.Stop(ctx); err != nil {
		log.Error("HTTP server shutdown", slog.Any("error", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
