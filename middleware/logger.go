package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &wrappedWriter{ResponseWriter: w, statusCode: 200}
			startTime := time.Now()
			next.ServeHTTP(writer, r)
			latency := time.Since(startTime).Microseconds()
			slog.Info(
				"incoming request",
				"path", r.URL.Path,
				"statusCode", writer.statusCode,
				"method", r.Method,
				//"userAgent", r.UserAgent(),
				"latencyMicro", latency,
			)
		})
	}
}
