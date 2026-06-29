package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.com/api/internals/logger"
	"example.com/api/internals/utils"
	"golang.org/x/time/rate"
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
	limiterMap := make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Context().Value(REQ_IP_KEY)
		castIP, ok := ip.(string)
		if !ok {
			http.Error(w, "Bad Request!", http.StatusBadRequest)
			return
		}

		mu.Lock()
		
		limiter, ok := limiterMap[castIP]
		if !ok {
			limiter = rate.NewLimiter(app.Config.RlLimit, app.Config.RlBurst)
			limiterMap[castIP] = limiter
		}

		mu.Unlock()

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(map[string]string{"message": "too many requests"}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		next.ServeHTTP(w, r)
	}
}

// GitHub webhooks request IP is not necessarily deterministic. until not verified, don't use this mw
func (app *Application) CheckAllowedDomainsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check req ip for now
		ip := r.Context().Value(REQ_IP_KEY)

		allowedDomainsFromEnv := utils.GetEnvString(utils.ALLOWED_DOMAINS)
		if allowedDomainsFromEnv == "" {
			http.Error(w, errors.New("internal server error").Error(), http.StatusInternalServerError)
		}

		allowedDomains := strings.Split(allowedDomainsFromEnv, ";")
		if !slices.Contains(allowedDomains, ip.(string)) {
			http.Error(w, "Not Allowed!", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// ValidateGitHubWebhookMiddleware validates HMAC signature coming from X-Hub-Signature-256 header,
// and event type in request body
func (app *Application) ValidateGitHubWebhookMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerSignature := r.Header.Get("X-Hub-Signature-256")
		if headerSignature == "" {
			http.Error(w, "Not Authorized!", http.StatusUnauthorized)
		}
		
		ghSecret := utils.GetEnvString(utils.GITHUB_SECRET)
		if ghSecret == "" {
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		}
		defer r.Body.Close()

		if !strings.HasPrefix(headerSignature, "sha256=") {
			http.Error(w, "Not Authorized!", http.StatusUnauthorized)
		}

		headerSignature = strings.TrimPrefix(headerSignature, "sha256=")

		mac := hmac.New(sha256.New, []byte(headerSignature))
		throwawayBody := body
		mac.Write(throwawayBody)
		expectedMAC := mac.Sum(nil)
		expectedSig := hex.EncodeToString(expectedMAC)

		if subtle.ConstantTimeCompare([]byte(headerSignature), []byte(expectedSig)) != 1 {
			http.Error(w, "Not Authorized!", http.StatusUnauthorized)
		}

		// TODO:
		// validate event type
		xGHEventHeader := r.Header.Get("x-github-event")
		fmt.Println("xGHEventHeader:", xGHEventHeader)

		reqBody := map[string]any{}
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		}
		fmt.Println("req. body:", reqBody)
		

		next.ServeHTTP(w, r)
	}
}