package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"

	"golang-url-shortner/internal/config"
	mwLogger "golang-url-shortner/internal/http-server/middleware/logger"
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

	router := chi.NewRouter()

	//making for easy finding lines in Grafana/Kibana via request-id (grep)
	router.Use(middleware.RequestID)
	//for getting real ip
	router.Use(middleware.RealIP)
	//logger from chi for more information
	router.Use(middleware.Logger)
	//custom logger from internal
	router.Use(mwLogger.New(log))
	//panic from handler + no drop application
	router.Use(middleware.Recoverer)
	//write urls for api
	router.Use(middleware.URLFormat)

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
