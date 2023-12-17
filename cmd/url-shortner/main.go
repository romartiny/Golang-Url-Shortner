package main

import (
	"log/slog"
	"os"

	"golang-url-shortner/internal/config"
	"golang-url-shortner/internal/lib/logger/sl"
	"golang-url-shortner/internal/storage/sqlite"
)

// environment constants
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//init config: cleanenv
	cfg := config.MustLoad()

	//init logger: sl
	log := setupLogger(cfg.Env)
	log.Info("starting golang-url-shortner", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1) //final error
	}

	_ = storage

	//init router: chi & chi render

	//run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	//switch env for different level debugging
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
