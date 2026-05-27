package main

import (
	"log"
	"net/http"
	"time"
)

type Application struct {
	Config Config
	DB any // not implemented
	StartedAt time.Time
	Runner Runner
}

type Config struct {
	Addr string
}

func NewApplication(address string) *Application{
	return &Application{
		Config: Config{
			Addr: address,
		},
		StartedAt: time.Now(),
		Runner: Runner{
			LastWorkflowSinceStart: LastWorkflowStat{
				Start: time.Time{},
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
		Addr: app.Config.Addr,
		Handler: middlewareStack(router),
	}
	
	return srv.ListenAndServe()
}