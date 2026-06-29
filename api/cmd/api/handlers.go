package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/api/internals/runner"
)

type AppStatus struct {
	IsAlive bool   `json:"is_alive"`
	Uptime  string `json:"uptime"`
	runner.LastWorkflowStat
}

type Log struct {
	Id            int64  `json:"id"`
	Timestamp     string `json:"timestamp"`
	MethodAndPath string `json:"method_and_path"`
	ResponseCode  int64  `json:"response_code"`
}

var logs = []Log{
	{
		Id:            1,
		Timestamp:     time.Now().UTC().String(),
		MethodAndPath: "GET /v1/health",
		ResponseCode:  200,
	},
	{
		Id:            2,
		Timestamp:     time.Now().UTC().String(),
		MethodAndPath: "GET /v1/stats",
		ResponseCode:  200,
	},
}

func (app *Application) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	appStatus := AppStatus{
		IsAlive: true,
		Uptime:  time.Since(app.StartedAt).Truncate(time.Second).String(),
		LastWorkflowStat: runner.LastWorkflowStat{
			Start:  app.Runner.LastWorkflowSinceStart.Start,
			Finish: app.Runner.LastWorkflowSinceStart.Finish,
		},
	}

	err := json.NewEncoder(w).Encode(appStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) getStatsPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := r.URL.Query().Get("limit")
	parsedLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil && limit != "" {
		http.Error(w, "Bad Request!", http.StatusBadRequest)
		return
	}

	page := r.URL.Query().Get("page")
	parsedPage, err := strconv.ParseInt(page, 10, 64)
	if err != nil && page != "" {
		http.Error(w, "Bad Request!", http.StatusBadRequest)
		return
	}

	log.Println("limit:", parsedLimit, "page:", parsedPage)

	// TODO: get last X logs, ignore limit and page for now
	err = json.NewEncoder(w).Encode(logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// workflow execution runs in a goroutine, handler responds if one is already running
func (app *Application) triggerCICDWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")

	if app.Runner.IsWorkflowRunning {
		if err := json.NewEncoder(w).Encode("{'workflow': 'running'}"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	app.Runner.IsWorkflowRunning = true
	app.Runner.LastWorkflowSinceStart.Start = time.Now().UTC()

	app.GlobalWG.Go(func() {
		if err := app.Runner.ExecutePipeline(context.Background(), app.ShutdownChan); err != nil {
			fmt.Printf("execute pipeline error: %v\n", err)
		}
	})

	if err := json.NewEncoder(w).Encode("{'workflow': 'initiating'}"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *Application) wildcardRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found!", http.StatusNotFound)
}
