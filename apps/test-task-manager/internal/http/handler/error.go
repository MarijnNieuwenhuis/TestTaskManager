package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type errorResponse struct {
	Error string `json:"error"`
}

func errorHandler(err error, code int, w http.ResponseWriter, logger *zap.SugaredLogger) {
	if err == nil {
		return
	}

	if logger != nil {
		if code >= 500 {
			logger.Errorw("internal server error", "error", err)
		} else {
			logger.Warnw("client error", "error", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(errorResponse{
		Error: err.Error(),
	})
}