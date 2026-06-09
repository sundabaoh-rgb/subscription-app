package http

import (
	"net/http"
	"subServ/internal/domain"
	"time"
)

func LoggingMiddleware(log domain.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(w, r)

		log.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}
