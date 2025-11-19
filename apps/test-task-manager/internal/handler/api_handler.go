// Package handler implements HTTP request handlers.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/btcdirect-api/test-task-manager/internal/service"
	"gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

// APIHandler handles JSON API requests.
type APIHandler struct {
	service *service.TaskService
}

// NewAPIHandler creates a new APIHandler.
func NewAPIHandler(service *service.TaskService) *APIHandler {
	return &APIHandler{service: service}
}

// GetTasks returns all tasks as JSON.
func (h *APIHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.GetAll()
	respondJSON(w, tasks, http.StatusOK)
}

// CreateTask creates a new task from JSON.
func (h *APIHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", "INVALID_INPUT", http.StatusBadRequest)
		return
	}

	task, err := h.service.Create(req.Title)
	if err != nil {
		if errors.Is(err, service.ErrEmptyTitle) || errors.Is(err, service.ErrTitleTooLong) {
			respondError(w, err.Error(), "INVALID_INPUT", http.StatusBadRequest)
			return
		}
		respondError(w, "Failed to create task", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	respondJSON(w, task, http.StatusCreated)
}

// ToggleTask toggles task completion status.
func (h *APIHandler) ToggleTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := h.service.Toggle(id)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			respondError(w, "Task not found", "NOT_FOUND", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to toggle task", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	respondJSON(w, task, http.StatusOK)
}

// DeleteTask deletes a task.
func (h *APIHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			respondError(w, "Task not found", "NOT_FOUND", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to delete task", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	respondJSON(w, MessageResponse{Message: "Task deleted successfully"}, http.StatusOK)
}
