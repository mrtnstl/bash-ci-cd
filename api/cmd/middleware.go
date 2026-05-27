package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

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
	return func(w http.ResponseWriter, r *http.Request){
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
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
	return func (w http.ResponseWriter, r *http.Request) {
		// TODO:
		next.ServeHTTP(w, r)
	}
}

func (app *Application) CheckAllowedDomainsMiddleware(next http.Handler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		ip, err := getIP(r)
		if err != nil {
			http.Error(w, "Can not determine client origin", http.StatusInternalServerError)
			return
		}
		
		allowedDomainsFromEnv := getEnvString("ALLOWED_DOMAINS")
		
		allowedDomains := strings.Split(allowedDomainsFromEnv, ";")
		if !slices.Contains(allowedDomains, ip) {
			http.Error(w, "Not Allowed!", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}