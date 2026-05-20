package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"example.com/api/internal"
)

type AppStatus struct {
	IsAlive bool `json:"is_live"`
	Uptime string `json:"uptime"`
	LastWorkflowStat
}

type Log struct {
	Id int64 `json:"id"`
	Timestamp string `json:"timestamp"`
	MethodAndPath string `json:"method_and_path"`
	ResponseCode int64 `json:"response_code"`
}

var logs = []Log{
	{
		Id: 1,
		Timestamp: time.Now().UTC().String(),
		MethodAndPath: "GET /v1/health",
		ResponseCode: 200,
	},
	{
		Id: 2,
		Timestamp: time.Now().UTC().String(),
		MethodAndPath: "GET /v1/stats",
		ResponseCode: 200,
	},
}

func (app *Application) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	appStatus := AppStatus{
		IsAlive: true,
		Uptime: time.Since(app.StartedAt).Truncate(time.Second).String(),
		LastWorkflowStat: LastWorkflowStat{
			Start: app.LastWorkflowSinceStart.Start,
			Finish: app.LastWorkflowSinceStart.Finish,
		},
	}

	err := json.NewEncoder(w).Encode(appStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (app *Application) getStatsPaginatedHandler(w http.ResponseWriter, r *http.Request){
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
func (app *Application) triggerCICDWorkflowHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if app.IsWorkflowRunning {
		err := json.NewEncoder(w).Encode("{'workflow': 'running'}")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusContinue)
		return
	}

	app.IsWorkflowRunning = true
	app.LastWorkflowSinceStart.Start = time.Now().UTC()

	var wg sync.WaitGroup
	errorChan := make(chan error, 1)

	wg.Go(func(){
		if err := internal.ExecutePipeline(r.Context()); err != nil {
			errorChan <- err		
		}
	})

	wg.Wait()
	close(errorChan)

	var errsFromGoroutine []error
	for err := range errorChan {
		errsFromGoroutine = append(errsFromGoroutine, err)
	}

	if len(errsFromGoroutine) != 0 {
		app.LastWorkflowSinceStart.Finish = time.Now().UTC()
		app.IsWorkflowRunning = false

		http.Error(w, errsFromGoroutine[0].Error(), http.StatusExpectationFailed)
		return
	}

	app.LastWorkflowSinceStart.Finish = time.Now().UTC()

	err := json.NewEncoder(w).Encode("{'workflow': 'finished'}")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	
	app.IsWorkflowRunning = false
}

func (app *Application) wildcardRouteHandler(w http.ResponseWriter, r *http.Request){
	http.Error(w, "Not Found!",http.StatusNotFound)
}