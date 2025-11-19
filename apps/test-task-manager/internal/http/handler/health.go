package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/btcdirect-api/test-task-manager/internal/app"
)

type configProvider interface {
	Config() app.Configuration
}

// HealthHandler returns a 200 OK status code.
func HealthHandler(provider configProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type output struct {
			Environment string `json:"environment"`
		}

		o := output{
			Environment: string(provider.Config().Environment),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(o)
	}
}

// ReadinessHandler returns a 200 OK status code if the database connection is alive.
// Otherwise, it returns a 503 Service Unavailable status code.
func ReadinessHandler(dbConn interface {
	IsAlive() bool
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type output struct {
			DatabaseHealthy bool `json:"databaseHealthy"`
		}

		o := output{
			DatabaseHealthy: dbConn != nil && dbConn.IsAlive(),
		}

		w.Header().Set("Content-Type", "application/json")
		defer json.NewEncoder(w).Encode(o)

		if !o.DatabaseHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
