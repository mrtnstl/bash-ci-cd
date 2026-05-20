package cmd

import (
	"log"
	"net/http"
	"time"
)

type LastWorkflowStat struct {
	Start time.Time `json:"last_wf_start"`
	Finish time.Time `json:"last_wf_finish"`
}

type Application struct {
	Config Config
	DB any // not implemented
	StartedAt time.Time
	LastWorkflowSinceStart LastWorkflowStat
	IsWorkflowRunning bool
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
		LastWorkflowSinceStart: LastWorkflowStat{
			Start: time.Time{},
			Finish: time.Time{},
		},
		IsWorkflowRunning: false,
	}
}

func (app *Application) Run() error {
	router := app.Mount()
	middlewareStack := app.MiddlewareStack(
		app.RequestLoggerMiddleware,
	)
	log.Printf("API started on port %s\n", app.Config.Addr)
	
	srv := &http.Server{
		Addr: app.Config.Addr,
		Handler: middlewareStack(router),
	}
	
	return srv.ListenAndServe()
}