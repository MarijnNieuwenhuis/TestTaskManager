package server

import (
	"gitlab.com/btcdirect-api/go-modules/http"
	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
)

type Server interface {
	Shutdown()
}

// Start Creates a new HTTP server, registers routes and starts it.
// Do not forget to call Shutdown() on the server when shutting down.
func Start(application *app.App) Server {
	s := http.CreateServer(application.Config().HTTPPort, application.Logger())

	registerRoutes(s.Router, application)

	s.Start()

	return s
}
