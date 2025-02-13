package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
        // выполняется только при старте приложения 1 раз
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

        // выполняется при каждом запросе
		fn := func(w http.ResponseWriter, r *http.Request) {
            // выполняется преред запросом
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
            // будет вуполнен после окончательной обработке запроса
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

            // передаем управление следующему handler в цепочке
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
