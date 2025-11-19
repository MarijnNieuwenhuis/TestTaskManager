package server

import (
	"gitlab.com/btcdirect-api/go-modules/http"
	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
	"gitlab.com/btcdirect-api/test-task-manager/internal/handler"
	"gitlab.com/btcdirect-api/test-task-manager/internal/service"
	"gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

type Server interface {
	Shutdown()
}

// Start Creates a new HTTP server, registers routes and starts it.
// Do not forget to call Shutdown() on the server when shutting down.
func Start(application *app.App) Server {
	s := http.CreateServer(application.Config().HTTPPort, application.Logger())

	// Initialize task manager components
	taskStore := store.NewTaskStore()
	taskService := service.NewTaskService(taskStore)
	pageHandler := handler.NewPageHandler(taskService)
	apiHandler := handler.NewAPIHandler(taskService)

	registerRoutes(s.Router, application, pageHandler, apiHandler)

	s.Start()

	return s
}
