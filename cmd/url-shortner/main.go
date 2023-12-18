package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang-url-shortner/internal/http-server/handlers/redirect"
	save "golang-url-shortner/internal/http-server/handlers/url"
	mwLogger "golang-url-shortner/internal/http-server/middleware/logger"
	"golang.org/x/exp/slog"
	"net/http"
	"os"

	"golang-url-shortner/internal/config"
	"golang-url-shortner/internal/lib/logger/handlers/slogpretty"
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
	log.Info(
		"starting golang-url-shortner",
		slog.String("env", cfg.Env),
		slog.String("version", cfg.Version),
	)
	log.Debug("debug messages are enabled")

	//init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1) //final error
	}

	_ = storage

	//init router: chi & chi render
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

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("golang-url-shortner", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
	})

	//routes
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")

	//run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	//switch env for different level debugging
	switch env {
	case envLocal:
		//init slogpretty
		log = setupPrettySlog()
		//for normal log
		//log = slog.New(
		//	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		//)
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

// prettier for logging
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
