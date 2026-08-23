package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	status              int
	size                int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	size, err := r.ResponseWriter.Write(b)
	r.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if r.status != 0 {
		return
	}

	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}

	return w.status
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.

func RequestLoggerMiddleware(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			}

			h.ServeHTTP(&lw, r)
			duration := time.Since(start)

			logger.Info("got incoming HTTP request",
				zap.String("Method", r.Method),
				zap.String("Path", r.URL.Path),
				zap.Duration("Duration", duration),
				zap.Int("Status", lw.Status()),
				zap.Int("Size", lw.size),
			)
		})
	}
}
