package api

import (
	"log"
	"net/http"
	"time"

	"example.com/api/internals/runner"
)

type Application struct {
	Config    Config
	DB        any // not implemented
	StartedAt time.Time
	Runner    runner.Runner
}

type Config struct {
	Addr string
}

func NewApplication(config Config) *Application {
	return &Application{
		Config: config,
		StartedAt: time.Now(),
		Runner: runner.Runner{
			LastWorkflowSinceStart: runner.LastWorkflowStat{
				Start:  time.Time{},
				Finish: time.Time{},
			},
			IsWorkflowRunning: false,
		},
	}
}

func (app *Application) Run() error {
	router := app.Mount()
	middlewareStack := app.MiddlewareStack(
		app.RequestLoggerMiddleware,
		app.CheckAllowedDomainsMiddleware,
		app.RateLimiterMiddleware,
		app.RequireHeaderSecretMiddleware,
	)
	log.Printf("API started on port %s\n", app.Config.Addr)

	srv := &http.Server{
		Addr:    app.Config.Addr,
		Handler: middlewareStack(router),
	}

	return srv.ListenAndServe()
}
