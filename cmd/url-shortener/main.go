package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
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

    // TODO: init router: chi, chi render

    // TODO run server
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



