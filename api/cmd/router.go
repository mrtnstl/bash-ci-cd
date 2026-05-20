package cmd

import (
	"net/http"
)

func (app *Application) Mount() *http.ServeMux{
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.wildcardRouteHandler)

	v1 := http.NewServeMux()

	v1.HandleFunc("GET /health", app.getHealthHandler)
	v1.HandleFunc("GET /stats", app.getStatsPaginatedHandler)
	v1.HandleFunc("POST /trigger", app.triggerCICDWorkflowHandler)
	v1.HandleFunc("/", app.wildcardRouteHandler)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	return mux
}