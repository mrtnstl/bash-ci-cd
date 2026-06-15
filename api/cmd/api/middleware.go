package api

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"example.com/api/internals/logger"
	"example.com/api/internals/utils"
)

const REQ_IP_KEY = "reqIP"

// resp writer wrapper for middleware that need response status code
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type Middleware func(http.Handler) http.HandlerFunc

func (app *Application) MiddlewareStack(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}

		return next.ServeHTTP
	}
}

func (app *Application) RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqIP, err := utils.GetIP(r)
		if err != nil || reqIP == "" {
			http.Error(w, "Can't determine client origin", http.StatusInternalServerError)
			return
		}

		writerContextWithIP := context.WithValue(r.Context(), REQ_IP_KEY, reqIP)
		r = r.WithContext(writerContextWithIP)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		statusAsStr := strconv.Itoa(wrapped.statusCode)
		reqDurationAsStr := time.Since(start).String()

		logger.Log(app.Config.AccessLogLocation, statusAsStr, r.Method, r.URL.Path, reqDurationAsStr, reqIP)
	}
}

func (app *Application) RequireHeaderSecretMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Not Authorized!", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *Application) RateLimiterMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		next.ServeHTTP(w, r)
	}
}

func (app *Application) CheckAllowedDomainsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Context().Value(REQ_IP_KEY)

		allowedDomainsFromEnv := utils.GetEnvString(utils.ALLOWED_DOMAINS)

		allowedDomains := strings.Split(allowedDomainsFromEnv, ";")
		if !slices.Contains(allowedDomains, ip.(string)) {
			http.Error(w, "Not Allowed!", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
