package server

import (
	"github.com/gorilla/mux"
	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
	"gitlab.com/btcdirect-api/test-task-manager/internal/http/handler"
)

// Registers all routes for the application.
func registerRoutes(r *mux.Router, app *app.App) {
	r.HandleFunc("/health", handler.HealthHandler(app)).Methods("GET")

	// TODO: Add your application-specific routes here
}
