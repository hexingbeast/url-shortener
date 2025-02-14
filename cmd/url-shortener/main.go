package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	deleteurl "url-shortener/internal/http-server/handlers/delete"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
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
    // ----------------------------------------------------
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
    // ----------------------------------------------------

    // TODO: init router: chi, chi render
    // create router
    router := chi.NewRouter()

    // connect to our router middleware
    // middleware it is handler, can be chain of handlers
    // last handler in chain, names request handler

    // middleware for adding requestId in our request(get it from chi dependency)
    router.Use(middleware.RequestID)
    // middleware for adding logger(get it from chi dependency, default logger)
    router.Use(middleware.Logger)
    // added custom middleware logger
    // можно использовать вместо дефолтного логгера chi
    router.Use(mwLogger.New(log))
    // если случилась паника внутри handler, не должны останавливать все приложение
    // из-за ошибки в одном handler, восстанавливаем от паники
    router.Use(middleware.Recoverer)
    // надо для того чтобы писать красивые url при подключению к нашему handler
    // привязан к пакету chi
    router.Use(middleware.URLFormat)

    // для авторизации добавляем новый роутер
    // в него помещаем handlers, для которых будет нужна авторизация 
    router.Route("/url", func(r chi.Router) {
        // BasicAuth идет в chi из коробки(претпологает оправку login и password в заголовке)
        r.Use(middleware.BasicAuth("url-shortener", map[string]string{
            cfg.HTTPServer.User: cfg.HTTPServer.Password,
        }))
        
        // чтобы заработала авторизация, надо перенести роутеры сюда
        // и меняем на внутренний роутер с "router" на "r"
        r.Post("/", save.New(log, storage))
        r.Delete("/{alias}", deleteurl.New(log, storage))
    })

    // добавляем handler для сохранения запроса
    // router.Post("/url", save.New(log, storage))
    router.Get("/{alias}", redirect.New(log, storage))
    // router.Delete("url/{alias}", deleteurl.New(log, storage))

    // TODO: run server
    log.Info("starting server", slog.String("address", cfg.Address))

    // создаем сам сервер
    srv := &http.Server{
        Addr: cfg.Address,
        // router также является handler, получается что это handler
        // внутри с нашими добавленными handler-ами
        Handler: router,
        ReadTimeout: cfg.HTTPServer.Timeout,
        WriteTimeout: cfg.HTTPServer.Timeout,
        IdleTimeout: cfg.HTTPServer.IdleTimeout,
    }

    // вызываем наш сервер, ListenAndServe() это блокирующая функция
    // она не пускает нашу программу дальше
    if err := srv.ListenAndServe(); err != nil {
        log.Error("failed to start server")
    }
    // если сюда программа дошла, то произошла ошибка и сервер остановлен
    log.Error("server stopped")
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

