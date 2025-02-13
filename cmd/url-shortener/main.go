package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
    envLocal = "local"
    envDev = "dev"
    envProd = "prod"
)

func main() {
    // TODO: init config: cleanenv
    cfg := config.MustLoad()
    // fmt.Println(cfg)

    // TODO: init logger: slog import from "log/slog"
    log := setupLogger(cfg.Env)
    log.Info("start url-shortener", slog.String("env", cfg.Env))
    // log.Debug("debug message are enable")

    // TODO: init storage: sqlite
    storage, err := sqlite.New(cfg.StoragePath)
    if err != nil {
        log.Error("failed to init storage", sl.Err(err))
        os.Exit(1)
    }
    // code for DEBUG
    // _, err1 := storage.SaveURL("google1.com", "google1")
    // if err1 != nil {
    //     log.Error("failed to save url", sl.Err(err1))
    //     os.Exit(1)
    // }
    //
    // url, err2 := storage.GetURL("google1")
    // if err2 != nil {
    //     log.Error("failed to get url", sl.Err(err2))
    // }
    // log.Info("get url ", slog.String("url", url))
    // log.Info("save url", slog.Int64("id", id))
    //
    // id2, err2 := storage.SaveURL("google.com", "google")
    // if err2 != nil {
    //     log.Error("failed to save url", sl.Err(err2))
    //     os.Exit(1)
    // }
    // log.Info("save url", slog.Int64("id", id2))

   _ = storage

    // TODO: init router: chi, chi render
    // create router
    router := chi.NewRouter()

    // connect to our router middleware
    // middleware it is handler, can be chain of handlers
    // last handler in chain, names request handler

    // middleware for adding requestId in our request(get it from chi dependency)
    router.Use(middleware.RequestID)
    // middleware for adding logger(get it from chi dependency)
    router.Use(middleware.Logger)

    // TODO: run server
}

// setup logger for different enviroments
func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

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
    default:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        )
    }

    return log
}



