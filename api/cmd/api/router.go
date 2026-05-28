package api

import (
	"net/http"
)

func (app *Application) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, app.Config.Static, "swagger.html")
	})
	mux.HandleFunc("GET /swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, app.Config.Static, "openapi.yaml")
	})

	v1 := http.NewServeMux()

	v1.HandleFunc("GET /health", app.getHealthHandler)
	v1.HandleFunc("GET /stats", app.getStatsPaginatedHandler)
	v1.HandleFunc("POST /trigger", app.triggerCICDWorkflowHandler)
	v1.HandleFunc("/", app.wildcardRouteHandler)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	mux.HandleFunc("/", app.wildcardRouteHandler)

	return mux
}
