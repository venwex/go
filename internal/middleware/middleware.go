package middleware

import (
	u "example/test/internal/utils"
	"log"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	const value = "secret12345"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")

		if key != value {
			u.RenderError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(message string) func(http.Handler) http.Handler  {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, message)

			next.ServeHTTP(w, r)
		})
	}
}
