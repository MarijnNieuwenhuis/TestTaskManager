package app

import (
	"time"

	"gitlab.com/btcdirect-api/go-modules/app"
	"go.uber.org/zap"
)

type App struct {
	config Configuration
	core   *app.App
}

// Initialize the application.
// This will also load the configuration.
func Initialize(c Configuration) *App {
	// In development mode, we set the shutdown timeout to 0 to allow for instant shutdowns.
	// In production, we set it to 30 seconds to allow for graceful shutdowns.
	shutdownTimeout := 30 * time.Second
	if c.Environment == Dev {
		shutdownTimeout = 0
	}

	core := app.Initialize(
		app.WithLoggerForLevel(c.LogLevel),
		app.WithShutdownTimeout(shutdownTimeout),
	)

	return &App{
		config: c,
		core:   &core,
	}
}

// Run the application and its services.
func (a *App) Run() {
	a.core.Run()
}

// Shutdown shuts down all services of the application.
func (a *App) Shutdown() {
	// No additional cleanup needed
}

// Config returns the application configuration.
func (a *App) Config() Configuration {
	return a.config
}

// Logger exposes the shared structured logger.
func (a *App) Logger() *zap.SugaredLogger {
	return a.core.Log
}
