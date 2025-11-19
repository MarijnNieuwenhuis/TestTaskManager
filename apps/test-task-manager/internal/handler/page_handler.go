package handler

import (
	"html/template"
	"net/http"

	"gitlab.com/btcdirect-api/test-task-manager/internal/model"
	"gitlab.com/btcdirect-api/test-task-manager/internal/service"
)

// PageHandler handles HTML page requests.
type PageHandler struct {
	service   *service.TaskService
	templates *template.Template
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(service *service.TaskService) *PageHandler {
	// Parse all templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	return &PageHandler{
		service:   service,
		templates: templates,
	}
}

// ServeTaskList renders the main task list page.
func (h *PageHandler) ServeTaskList(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.GetAll()

	data := struct {
		Tasks []model.Task
	}{
		Tasks: tasks,
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
