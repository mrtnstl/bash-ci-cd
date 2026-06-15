package api

import (
	"embed"
	"log"
	"net/http"
	"sync"
	"time"

	"example.com/api/internals/runner"
)

type Application struct {
	Config       Config
	DB           any // not implemented
	StartedAt    time.Time
	Runner       runner.Runner
	GlobalWG     *sync.WaitGroup
	ShutdownChan *chan bool
}

type Config struct {
	Addr              string
	Static            embed.FS
	AccessLogLocation string
}

func NewApplication(config Config, wg *sync.WaitGroup, shutdownChan *chan bool) *Application {
	return &Application{
		Config:    config,
		StartedAt: time.Now(),
		Runner: runner.Runner{
			LastWorkflowSinceStart: runner.LastWorkflowStat{
				Start:  time.Time{},
				Finish: time.Time{},
			},
			IsWorkflowRunning: false,
		},
		GlobalWG:     wg,
		ShutdownChan: shutdownChan,
	}
}

func (app *Application) CreateServer() *http.Server {
	router := app.Mount()
	middlewareStack := app.MiddlewareStack(
		app.RequestLoggerMiddleware,
		app.CheckAllowedDomainsMiddleware,
		app.RateLimiterMiddleware,
		app.RequireHeaderSecretMiddleware,
	)

	srv := &http.Server{
		Addr:    app.Config.Addr,
		Handler: middlewareStack(router),
	}

	return srv
}

func (app *Application) Run(srv *http.Server) error {
	log.Printf("API started on port %s\n", app.Config.Addr)

	return srv.ListenAndServe()
}
