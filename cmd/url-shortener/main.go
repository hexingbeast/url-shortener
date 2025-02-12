package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
    envLocal = "local"
    envDev = "dev"
    envProd = "prod"
)

func main() {
    // TODO: init config: cleanenv
    cfg := config.MustLoad()
    fmt.Println(cfg)

    // TODO: init logger: slog import from "log/slog"
    log := setupLogger(cfg.Env)
    log.Info("start url-shortener", slog.String("env", cfg.Env))
    log.Debug("debug message are enable")

    // TODO: init storage: sqlite
    storage, err := sqlite.New(cfg.StoragePath)
    if err != nil {
        log.Error("failed to init storage", sl.Err(err))
        os.Exit(1)
    }
    
    // id, err := storage.SaveURL("google.com", "google")
    // if err != nil {
    //     log.Error("failed to save url", sl.Err(err))
    //     os.Exit(1)
    // }
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

    // TODO: run server
}

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



