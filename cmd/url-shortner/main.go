package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"url-shortner/internal/config"
	"url-shortner/internal/http-server/middleware/logger"
	"url-shortner/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.NewConfig()

	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug message")
	fmt.Println()

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", err)
		return
	}

	_ = storage

	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server

	http.ListenAndServe(cfg.HTTPServer.Address, router)
}

const (
	envLocal = "local"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		// TODO: setup logger for dev
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic(fmt.Sprintf("unknown env %s", env))
	}

	return log
}
