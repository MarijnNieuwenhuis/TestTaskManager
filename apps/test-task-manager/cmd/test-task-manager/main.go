package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
	"gitlab.com/btcdirect-api/test-task-manager/internal/http/server"
)

func main() {
	c := app.Configuration{}

	var env string
	flag.StringVar(&env, "env", getenv("APP_ENV", "dev"), "Environment")

	var err error
	c.Environment, err = getEnvironment(env)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&c.LogLevel, "loglevel", getenv("LOG_LEVEL", "info"), "Log output level")
	flag.StringVar(&c.HTTPPort, "port", getenv("HTTP_PORT", "8080"), "HTTP port")

	flag.Parse()

	application := app.Initialize(c)

	run(application)
}

// Run the application daemon.
func run(application *app.App) {
	application.Logger().Info("Starting application")

	server := server.Start(application)
	application.Run()

	application.Logger().Info("Shutting down application")

	application.Shutdown()
	server.Shutdown()

	os.Exit(0)
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getEnvironment(input string) (app.Environment, error) {
	switch input {
	case "dev":
		return app.Dev, nil
	case "stage":
		return app.Stage, nil
	case "acc":
		return app.Acc, nil
	case "sandbox":
		return app.Sandbox, nil
	case "prod":
		return app.Prod, nil
	default:
		return "", fmt.Errorf("invalid environment: %s", input)
	}
}
