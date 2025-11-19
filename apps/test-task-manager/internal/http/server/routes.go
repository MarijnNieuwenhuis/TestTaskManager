package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
	"gitlab.com/btcdirect-api/test-task-manager/internal/handler"
	oldhandler "gitlab.com/btcdirect-api/test-task-manager/internal/http/handler"
)

// Registers all routes for the application.
func registerRoutes(r *mux.Router, app *app.App, pageHandler *handler.PageHandler, apiHandler *handler.APIHandler) {
	// Health endpoint
	r.HandleFunc("/health", oldhandler.HealthHandler(app)).Methods("GET")

	// Static files
	staticDir := http.Dir("static")
	staticHandler := http.StripPrefix("/static/", http.FileServer(staticDir))
	r.PathPrefix("/static/").Handler(staticHandler)

	// Page routes (HTML)
	r.HandleFunc("/", pageHandler.ServeTaskList).Methods("GET")

	// API routes (JSON)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", apiHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", apiHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}/toggle", apiHandler.ToggleTask).Methods("PATCH")
	api.HandleFunc("/tasks/{id}", apiHandler.DeleteTask).Methods("DELETE")
}
