# Feature: Simple Task Manager - Implementation Plan

**Status**: ✅ Complete
**Created**: 2025-11-19
**Completed**: 2025-11-19
**Go Version**: 1.25
**Architecture**: [ARCHITECTURE.md](./ARCHITECTURE.md)
**Requirements**: [FEATURE.md](./FEATURE.md)

## Overview

This implementation plan details the step-by-step tasks required to build a lightweight task management web application using Go (backend) and Stimulus.js (frontend). The application will use in-memory storage with thread-safe access via `sync.RWMutex`, serve HTML templates, and provide a REST API for dynamic AJAX operations.

## Architecture Summary

**Pattern**: Monolithic Web Application
**Components**:
- **Data Layer**: Task model + In-memory store (thread-safe with RWMutex)
- **Business Logic Layer**: TaskService (validation and business rules)
- **HTTP Layer**: Page handler (HTML) + API handler (JSON) + Router
- **Template Layer**: Go html/template with Bootstrap 5.3
- **Frontend Layer**: Stimulus.js controllers for dynamic interactions

## Implementation Phases

### Phase 1: Foundation & Project Setup ✅
- [x] Task 1: Clean up bootstrap template (remove unused components)
- [x] Task 2: Create directory structure for task manager
- [x] Task 3: Define Task model
- [x] Task 4: Implement in-memory TaskStore with thread safety

### Phase 2: Backend Core Implementation ✅
- [x] Task 5: Implement TaskService with validation
- [x] Task 6: Implement API handlers (JSON responses)
- [x] Task 7: Implement page handler (HTML rendering)
- [x] Task 8: Configure router with all endpoints

### Phase 3: Frontend Implementation ✅
- [x] Task 9: Create HTML templates with Bootstrap
- [x] Task 10: Add Stimulus.js and create tasks controller
- [x] Task 11: Set up static file serving

### Phase 4: Integration & Polish ✅
- [x] Task 12: Wire up main.go with dependency injection
- [x] Task 13: Add error handling and logging
- [x] Task 14: Update Makefile and configuration
- [x] Task 15: Manual testing and bug fixes

## Detailed Task Breakdown

---

### Task 1: Clean Up Bootstrap Template

**Phase**: Foundation
**Dependencies**: None
**Location**: `apps/test-task-manager/`
**Estimated Effort**: Small

#### Description
Remove unnecessary components from the bootstrap Go project template that aren't needed for the task manager (database, Pub/Sub messaging, Sentry, etc.).

#### Acceptance Criteria
- [ ] Database-related code removed from `internal/app/` and `internal/db/`
- [ ] Pub/Sub messenger code removed from `internal/messenger/`
- [ ] Sentry integration removed
- [ ] `.env` file updated (remove DATABASE_URL, SENTRY_DSN, PUBSUB_* variables)
- [ ] Unused dependencies removed from `go.mod`
- [ ] Project still compiles after cleanup

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Simplifications from Bootstrap Template"
- Keep: Gorilla Mux, logging infrastructure, basic app structure
- Remove: Database (sqlx, migrations), Cloud SQL, Pub/Sub, Sentry
- Update go.mod: Run `go mod tidy` after removing imports

#### Files to Create/Modify
- `internal/app/app.go` - Simplify app initialization (remove DB, Pub/Sub)
- Remove `internal/db/` directory entirely
- Remove `internal/messenger/` directory entirely
- `.env` - Remove unused environment variables
- `go.mod` - Remove unused dependencies via `go mod tidy`

#### Go Best Practices
- Reference: `.claude/refs/go/best-practices.md` - Keep it Simple principle
- Don't remove code that might be needed; only remove what's clearly unused

---

### Task 2: Create Directory Structure

**Phase**: Foundation
**Dependencies**: Task 1
**Location**: `apps/test-task-manager/internal/`
**Estimated Effort**: Small

#### Description
Create the directory structure for the task manager following the architecture design.

#### Acceptance Criteria
- [ ] `internal/model/` directory created
- [ ] `internal/store/` directory created
- [ ] `internal/service/` directory created
- [ ] `internal/handler/` directory exists (may already exist from bootstrap)
- [ ] `internal/server/` directory exists (may already exist from bootstrap)
- [ ] `templates/` directory created in project root
- [ ] `templates/partials/` directory created
- [ ] `static/js/` directory created
- [ ] `static/js/controllers/` directory created
- [ ] `static/css/` directory created (optional, for custom styles)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Directory Structure"
- Follow Go convention: `internal/` for private packages
- Templates and static files at project root level

#### Files to Create/Modify
- Create directory structure with `mkdir -p` commands

#### Go Best Practices
- Reference: `.claude/refs/go/best-practices.md` - Project Structure
- Use `internal/` to prevent external imports

---

### Task 3: Define Task Model

**Phase**: Foundation
**Dependencies**: Task 2
**Location**: `internal/model/task.go`
**Estimated Effort**: Small

#### Description
Create the Task struct that represents a task entity with proper JSON tags for API serialization.

#### Acceptance Criteria
- [ ] Task struct defined with ID, Title, Completed, CreatedAt fields
- [ ] JSON tags added for API serialization
- [ ] Proper field types (ID as string, CreatedAt as time.Time)
- [ ] Package documentation added
- [ ] File compiles without errors

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Task Model"
- Use `string` for ID (will be generated as sequential integers converted to string)
- Use `time.Time` for CreatedAt (RFC3339 format in JSON)
- Add JSON tags in camelCase for frontend compatibility

#### Files to Create/Modify
- `internal/model/task.go` - Create Task struct

```go
package model

import "time"

// Task represents a single task item in the task manager.
type Task struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"createdAt"`
}
```

#### Go Best Practices
- Reference: `.claude/refs/go/idiomatic-go.md` - Struct Tags
- Use camelCase in JSON tags (JavaScript convention)
- Add godoc comments for exported types

---

### Task 4: Implement In-Memory TaskStore

**Phase**: Foundation
**Dependencies**: Task 3
**Location**: `internal/store/task_store.go`
**Estimated Effort**: Medium

#### Description
Implement thread-safe in-memory task storage using `sync.RWMutex` for concurrent access protection.

#### Acceptance Criteria
- [ ] TaskStore struct defined with tasks slice, nextID counter, and RWMutex
- [ ] NewTaskStore() constructor implemented
- [ ] GetAll() method with read lock
- [ ] GetByID() method with read lock
- [ ] Create() method with write lock
- [ ] Toggle() method with write lock
- [ ] Delete() method with write lock
- [ ] Proper use of defer for unlock
- [ ] Errors defined for "not found" cases
- [ ] ID generation as sequential integers converted to strings

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Task Store (In-Memory)"
- Reference: `.claude/refs/go/concurrency-patterns.md` - Mutex Patterns
- Use `sync.RWMutex`: RLock for reads, Lock for writes
- Always `defer unlock()` immediately after lock to prevent deadlocks
- Use `strconv.Itoa()` to convert integer IDs to strings
- GetAll() should return a copy of the slice, not a reference

#### Files to Create/Modify
- `internal/store/task_store.go` - Complete implementation
- `internal/store/errors.go` - Define ErrTaskNotFound

```go
package store

import (
    "errors"
    "fmt"
    "strconv"
    "sync"
    "time"

    "gitlab.com/btcdirect-api/test-task-manager/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

// TaskStore provides thread-safe in-memory task storage
type TaskStore struct {
    tasks  []model.Task
    nextID int
    mu     sync.RWMutex
}

// NewTaskStore creates a new TaskStore
func NewTaskStore() *TaskStore {
    return &TaskStore{
        tasks:  make([]model.Task, 0),
        nextID: 1,
    }
}

// GetAll returns all tasks
func (s *TaskStore) GetAll() []model.Task {
    s.mu.RLock()
    defer s.mu.RUnlock()

    // Return a copy to prevent external modification
    tasksCopy := make([]model.Task, len(s.tasks))
    copy(tasksCopy, s.tasks)
    return tasksCopy
}

// Create adds a new task
func (s *TaskStore) Create(title string) model.Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    task := model.Task{
        ID:        strconv.Itoa(s.nextID),
        Title:     title,
        Completed: false,
        CreatedAt: time.Now(),
    }

    s.tasks = append(s.tasks, task)
    s.nextID++

    return task
}

// GetByID returns a task by ID
func (s *TaskStore) GetByID(id string) (model.Task, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, task := range s.tasks {
        if task.ID == id {
            return task, nil
        }
    }

    return model.Task{}, ErrTaskNotFound
}

// Toggle changes completion status
func (s *TaskStore) Toggle(id string) (model.Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for i := range s.tasks {
        if s.tasks[i].ID == id {
            s.tasks[i].Completed = !s.tasks[i].Completed
            return s.tasks[i], nil
        }
    }

    return model.Task{}, ErrTaskNotFound
}

// Delete removes a task
func (s *TaskStore) Delete(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    for i, task := range s.tasks {
        if task.ID == id {
            // Remove task from slice
            s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
            return nil
        }
    }

    return ErrTaskNotFound
}
```

#### Go Best Practices
- Reference: `.claude/refs/go/concurrency-patterns.md` - Mutex Usage
- Always defer unlock immediately after lock
- Return copies of slices to prevent race conditions
- Use range over slice indices when modifying during iteration

---

### Task 5: Implement TaskService

**Phase**: Backend Core
**Dependencies**: Task 4
**Location**: `internal/service/task_service.go`
**Estimated Effort**: Medium

#### Description
Implement business logic layer with input validation and error handling.

#### Acceptance Criteria
- [ ] TaskService struct defined with store dependency
- [ ] NewTaskService() constructor implemented
- [ ] GetAll() method delegates to store
- [ ] Create() method with validation (title not empty, max 255 chars)
- [ ] Toggle() method delegates to store
- [ ] Delete() method delegates to store
- [ ] Custom errors for validation failures
- [ ] Title trimming before validation
- [ ] Proper error wrapping with context

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Task Service"
- Reference: `.claude/refs/go/error-handling.md` - Error Wrapping
- Validation rules from FEATURE.md:
  - Title must not be empty after trimming
  - Title must not exceed 255 characters
- Use `strings.TrimSpace()` for title trimming
- Wrap store errors with additional context

#### Files to Create/Modify
- `internal/service/task_service.go` - Complete implementation
- `internal/service/errors.go` - Define validation errors

```go
package service

import (
    "errors"
    "fmt"
    "strings"

    "gitlab.com/btcdirect-api/test-task-manager/internal/model"
    "gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

var (
    ErrEmptyTitle      = errors.New("task title cannot be empty")
    ErrTitleTooLong    = errors.New("task title cannot exceed 255 characters")
)

// TaskService handles business logic for tasks
type TaskService struct {
    store *store.TaskStore
}

// NewTaskService creates a new TaskService
func NewTaskService(store *store.TaskStore) *TaskService {
    return &TaskService{store: store}
}

// GetAll retrieves all tasks
func (s *TaskService) GetAll() []model.Task {
    return s.store.GetAll()
}

// Create creates a new task with validation
func (s *TaskService) Create(title string) (model.Task, error) {
    // Trim whitespace
    title = strings.TrimSpace(title)

    // Validate title
    if title == "" {
        return model.Task{}, ErrEmptyTitle
    }

    if len(title) > 255 {
        return model.Task{}, ErrTitleTooLong
    }

    // Create task
    task := s.store.Create(title)
    return task, nil
}

// Toggle toggles task completion status
func (s *TaskService) Toggle(id string) (model.Task, error) {
    task, err := s.store.Toggle(id)
    if err != nil {
        return model.Task{}, fmt.Errorf("failed to toggle task: %w", err)
    }
    return task, nil
}

// Delete removes a task
func (s *TaskService) Delete(id string) error {
    if err := s.store.Delete(id); err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    return nil
}
```

#### Go Best Practices
- Reference: `.claude/refs/go/error-handling.md` - Sentinel Errors
- Use `fmt.Errorf` with `%w` for error wrapping
- Define custom errors as package-level variables

---

### Task 6: Implement API Handlers

**Phase**: Backend Core
**Dependencies**: Task 5
**Location**: `internal/handler/api_handler.go`
**Estimated Effort**: Large

#### Description
Implement JSON API handlers for CRUD operations on tasks.

#### Acceptance Criteria
- [ ] APIHandler struct defined with service dependency
- [ ] NewAPIHandler() constructor implemented
- [ ] GetTasks() endpoint (GET /api/tasks) returns JSON array
- [ ] CreateTask() endpoint (POST /api/tasks) accepts and returns JSON
- [ ] ToggleTask() endpoint (PATCH /api/tasks/{id}/toggle) returns JSON
- [ ] DeleteTask() endpoint (DELETE /api/tasks/{id}) returns JSON
- [ ] Proper HTTP status codes (200, 201, 400, 404, 500)
- [ ] ErrorResponse struct for consistent error formatting
- [ ] Content-Type headers set correctly
- [ ] Request body parsing and validation
- [ ] Proper error responses with error codes

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "API Handler" and "API Specification"
- Reference: `.claude/refs/go/error-handling.md` - HTTP Error Responses
- Use Gorilla Mux `mux.Vars(r)` to extract URL parameters
- Set Content-Type: application/json for all responses
- Return proper status codes as per API specification
- Use json.NewDecoder for request bodies, json.NewEncoder for responses

#### Files to Create/Modify
- `internal/handler/api_handler.go` - Complete implementation
- `internal/handler/response.go` - ErrorResponse and helper functions

```go
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/gorilla/mux"
    "gitlab.com/btcdirect-api/test-task-manager/internal/service"
    "gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}

// Success response for delete
type MessageResponse struct {
    Message string `json:"message"`
}

// APIHandler handles JSON API requests
type APIHandler struct {
    service *service.TaskService
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(service *service.TaskService) *APIHandler {
    return &APIHandler{service: service}
}

// GetTasks returns all tasks as JSON
func (h *APIHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
    tasks := h.service.GetAll()

    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

// CreateTask creates a new task from JSON
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

// ToggleTask toggles task completion status
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

// DeleteTask deletes a task
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(MessageResponse{Message: "Task deleted successfully"})
}

// Helper function to send error responses
func respondError(w http.ResponseWriter, message, code string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}
```

#### Go Best Practices
- Reference: `.claude/refs/go/error-handling.md` - Error Checking with errors.Is
- Use `errors.Is()` to check for specific errors
- Set headers before writing status code
- Use `json.NewEncoder` for streaming writes

---

### Task 7: Implement Page Handler

**Phase**: Backend Core
**Dependencies**: Task 5
**Location**: `internal/handler/page_handler.go`
**Estimated Effort**: Medium

#### Description
Implement HTML page handler that renders the main task list page using Go templates.

#### Acceptance Criteria
- [ ] PageHandler struct defined with service and templates
- [ ] NewPageHandler() constructor loads templates
- [ ] ServeTaskList() method renders index.html with task data
- [ ] Template parsing uses `template.ParseGlob()`
- [ ] Proper error handling for template rendering
- [ ] Template data structure defined
- [ ] HTTP 500 error on template failure

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Page Handler"
- Use `html/template` (not `text/template`) for auto-escaping
- Parse templates in constructor, not on every request
- Pass tasks to template via anonymous struct

#### Files to Create/Modify
- `internal/handler/page_handler.go` - Complete implementation

```go
package handler

import (
    "html/template"
    "net/http"

    "gitlab.com/btcdirect-api/test-task-manager/internal/service"
)

// PageHandler handles HTML page requests
type PageHandler struct {
    service   *service.TaskService
    templates *template.Template
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(service *service.TaskService) *PageHandler {
    // Parse all templates
    templates := template.Must(template.ParseGlob("templates/*.html"))

    return &PageHandler{
        service:   service,
        templates: templates,
    }
}

// ServeTaskList renders the main task list page
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
```

#### Go Best Practices
- Reference: `.claude/refs/go/best-practices.md` - Template Usage
- Use `html/template` for XSS protection via auto-escaping
- Parse templates once in constructor (not per request)
- Use `template.Must()` to panic on parse errors at startup

---

### Task 8: Configure Router

**Phase**: Backend Core
**Dependencies**: Task 6, Task 7
**Location**: `internal/server/routes.go`
**Estimated Effort**: Small

#### Description
Configure Gorilla Mux router with all HTTP endpoints and static file serving.

#### Acceptance Criteria
- [ ] NewRouter() function created accepting page and API handlers
- [ ] Static file serving configured for `/static/` path
- [ ] GET / route mapped to page handler
- [ ] GET /api/tasks route mapped to API handler
- [ ] POST /api/tasks route mapped to API handler
- [ ] PATCH /api/tasks/{id}/toggle route mapped to API handler
- [ ] DELETE /api/tasks/{id} route mapped to API handler
- [ ] Proper HTTP methods specified with `.Methods()`

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Router Configuration"
- Use `PathPrefix` for API routes grouping
- Use `http.FileServer` for static files
- Use `http.StripPrefix` to remove `/static/` from file paths

#### Files to Create/Modify
- `internal/server/routes.go` - Complete router configuration

```go
package server

import (
    "net/http"

    "github.com/gorilla/mux"
    "gitlab.com/btcdirect-api/test-task-manager/internal/handler"
)

// NewRouter creates and configures the application router
func NewRouter(pageHandler *handler.PageHandler, apiHandler *handler.APIHandler) *mux.Router {
    r := mux.NewRouter()

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

    return r
}
```

#### Go Best Practices
- Reference: Gorilla Mux documentation
- Use subrouters for API versioning/grouping
- Specify HTTP methods explicitly with `.Methods()`

---

### Task 9: Create HTML Templates

**Phase**: Frontend
**Dependencies**: None (can be done in parallel with backend)
**Location**: `templates/`
**Estimated Effort**: Large

#### Description
Create Go HTML templates with Bootstrap 5.3 styling for the task manager UI.

#### Acceptance Criteria
- [ ] base.html created with HTML5 boilerplate and Bootstrap CDN
- [ ] index.html created extending base with task list layout
- [ ] partials/task-form.html created with input form
- [ ] partials/task-item.html created for individual task display
- [ ] Bootstrap classes applied for responsive design
- [ ] Stimulus data attributes added (data-controller, data-action, data-target)
- [ ] Proper template composition with {{define}} and {{template}}
- [ ] XSS-safe (Go templates auto-escape)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Template Layer"
- Use Bootstrap 5.3 from CDN (jsDelivr or official CDN)
- Use Stimulus.js from npm CDN (unpkg.com or esm.sh)
- Add data attributes for Stimulus: `data-controller="tasks"`
- Follow Bootstrap conventions for forms, buttons, lists

#### Files to Create/Modify
- `templates/base.html` - Base layout
- `templates/index.html` - Main page content
- `templates/partials/task-form.html` - Task creation form
- `templates/partials/task-item.html` - Task item display

See full template code in ARCHITECTURE.md "Template Layer" section.

#### Go Best Practices
- Use `{{define}}` and `{{template}}` for composition
- Go templates auto-escape HTML (XSS protection)
- Keep templates simple; complex logic belongs in Go code

---

### Task 10: Implement Stimulus.js Controller

**Phase**: Frontend
**Dependencies**: None (can be done in parallel)
**Location**: `static/js/`
**Estimated Effort**: Large

#### Description
Implement Stimulus.js controller for dynamic task interactions (create, toggle, delete) via AJAX.

#### Acceptance Criteria
- [ ] app.js created to initialize Stimulus application
- [ ] controllers/tasks_controller.js created extending Controller
- [ ] Targets defined (list, form, input)
- [ ] create() method handles form submission
- [ ] toggle() method handles completion toggle
- [ ] delete() method handles task deletion with confirmation
- [ ] Helper methods for DOM manipulation (add, update, remove)
- [ ] Proper async/await usage for fetch calls
- [ ] Error handling for failed API calls
- [ ] ES6 module syntax (import/export)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Frontend Layer (Stimulus.js)"
- Use `fetch()` API for AJAX requests
- Use ES6 async/await for cleaner async code
- Handle HTTP errors (check `response.ok`)
- Show user feedback for errors (could use alert or console.log)

#### Files to Create/Modify
- `static/js/app.js` - Stimulus application setup
- `static/js/controllers/tasks_controller.js` - Tasks controller

See full JavaScript code in ARCHITECTURE.md "Frontend Layer" section.

#### Go Best Practices
Not applicable (JavaScript, not Go)

---

### Task 11: Set Up Static File Serving

**Phase**: Frontend
**Dependencies**: Task 8
**Location**: `static/`
**Estimated Effort**: Small

#### Description
Ensure static directory exists and files can be served correctly.

#### Acceptance Criteria
- [ ] `static/` directory created
- [ ] `static/js/` subdirectory created
- [ ] `static/js/controllers/` subdirectory created
- [ ] `static/css/` subdirectory created (for custom CSS if needed)
- [ ] Static file serving tested (files accessible via /static/ URL)
- [ ] Stimulus.js loaded from CDN or npm package

#### Implementation Notes
- Static files already configured in router (Task 8)
- Stimulus.js can be loaded from CDN (https://unpkg.com/@hotwired/stimulus)
- Or installed via npm and bundled (simpler: use CDN for testing)

#### Files to Create/Modify
- Create directory structure
- Optionally: `static/css/custom.css` for additional styles

---

### Task 12: Wire Up Main Application

**Phase**: Integration
**Dependencies**: Task 4, Task 5, Task 6, Task 7, Task 8
**Location**: `cmd/test-task-manager/main.go`
**Estimated Effort**: Small

#### Description
Update main.go to initialize all components with proper dependency injection and start the HTTP server.

#### Acceptance Criteria
- [ ] Database and Pub/Sub initialization removed
- [ ] TaskStore initialized
- [ ] TaskService initialized with store
- [ ] PageHandler initialized with service
- [ ] APIHandler initialized with service
- [ ] Router initialized with handlers
- [ ] HTTP server started on configurable port
- [ ] Startup logging added
- [ ] HTTP_PORT environment variable supported (default 8080)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Server Initialization"
- Use constructor functions for dependency injection
- Read PORT from environment or default to 8080
- Use `log.Printf` for startup message
- Use `log.Fatal` for server errors

#### Files to Create/Modify
- `cmd/test-task-manager/main.go` - Update main function

```go
package main

import (
    "log"
    "net/http"
    "os"

    "gitlab.com/btcdirect-api/test-task-manager/internal/handler"
    "gitlab.com/btcdirect-api/test-task-manager/internal/server"
    "gitlab.com/btcdirect-api/test-task-manager/internal/service"
    "gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

func main() {
    // Initialize store
    taskStore := store.NewTaskStore()

    // Initialize service
    taskService := service.NewTaskService(taskStore)

    // Initialize handlers
    pageHandler := handler.NewPageHandler(taskService)
    apiHandler := handler.NewAPIHandler(taskService)

    // Initialize router
    router := server.NewRouter(pageHandler, apiHandler)

    // Get port from environment or use default
    port := os.Getenv("HTTP_PORT")
    if port == "" {
        port = "8080"
    }

    // Start server
    addr := ":" + port
    log.Printf("Starting Task Manager server on %s", addr)
    log.Fatal(http.ListenAndServe(addr, router))
}
```

#### Go Best Practices
- Reference: `.claude/refs/go/best-practices.md` - Dependency Injection
- Use constructor functions (New*) for initialization
- Keep main() simple: wire dependencies and start server
- Use log.Fatal for startup errors (exits with non-zero code)

---

### Task 13: Add Error Handling and Logging

**Phase**: Integration
**Dependencies**: Task 12
**Location**: `internal/handler/`, `internal/service/`
**Estimated Effort**: Medium

#### Description
Add comprehensive error handling and request logging throughout the application.

#### Acceptance Criteria
- [ ] All API handlers log requests (method, path, status)
- [ ] All errors logged with ERROR level
- [ ] Validation errors logged at INFO level
- [ ] Task operations logged at DEBUG level (optional)
- [ ] HTTP middleware for request logging (optional but recommended)
- [ ] Consistent error response format across all endpoints
- [ ] No panics (all errors handled gracefully)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Logging Strategy"
- Use Go's `log` package (simple) or existing logger from bootstrap
- Log format: `[timestamp] [level] message`
- Consider middleware for automatic request logging

#### Files to Create/Modify
- `internal/handler/api_handler.go` - Add logging to all methods
- `internal/handler/page_handler.go` - Add logging
- `internal/service/task_service.go` - Add operation logging
- Optional: `internal/server/middleware.go` - Request logging middleware

#### Go Best Practices
- Reference: `.claude/refs/go/error-handling.md` - Logging Errors
- Log errors before returning them
- Include context in log messages (task ID, operation, etc.)
- Don't log and return the same error (choose one)

---

### Task 14: Update Makefile and Configuration

**Phase**: Integration
**Dependencies**: Task 12
**Location**: Root directory
**Estimated Effort**: Small

#### Description
Update Makefile and .env to reflect the simplified task manager configuration.

#### Acceptance Criteria
- [ ] Makefile CMD updated to run test-task-manager
- [ ] Database migration targets removed from Makefile
- [ ] .env file cleaned up (removed DB, Pub/Sub, Sentry variables)
- [ ] HTTP_PORT added to .env (default: 8080)
- [ ] LOG_LEVEL added to .env (default: debug)
- [ ] README.md updated with task manager description

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - "Build Configuration"
- Remove migrate, migrate-down targets
- Keep run, build, test targets
- Update CMD to: `go run ./cmd/test-task-manager/main.go`

#### Files to Create/Modify
- `Makefile` - Update run command and remove database targets
- `.env` - Simplify to only HTTP_PORT and LOG_LEVEL
- `README.md` - Update with task manager description

#### Go Best Practices
- Keep Makefile simple and focused
- Use environment variables for configuration
- Don't commit sensitive data to .env

---

### Task 15: Manual Testing and Bug Fixes

**Phase**: Integration
**Dependencies**: All previous tasks
**Location**: N/A
**Estimated Effort**: Medium

#### Description
Manually test all functionality in the browser and fix any bugs discovered.

#### Acceptance Criteria
- [ ] Application starts successfully with `make run`
- [ ] Home page (/) loads and displays empty state
- [ ] Create task: Submit form adds task to list without page reload
- [ ] Create task: Empty title shows validation error
- [ ] Create task: Long title (>255 chars) shows validation error
- [ ] Toggle task: Clicking checkbox toggles completion status
- [ ] Toggle task: Visual styling updates (strikethrough, color change)
- [ ] Delete task: Confirmation prompt shown
- [ ] Delete task: Confirming removes task from list
- [ ] Delete task: Canceling keeps task in list
- [ ] All interactions work without page reload (AJAX)
- [ ] UI is responsive on desktop and mobile viewports
- [ ] No console errors in browser developer tools
- [ ] No Go errors in server logs during normal operation

#### Implementation Notes
- Test in multiple browsers (Chrome, Firefox, Safari if available)
- Test on different screen sizes (desktop, tablet, mobile)
- Use browser developer tools to inspect network requests
- Check server logs for errors during testing

#### Testing Checklist
From FEATURE.md "Success Criteria":
- [ ] All CRUD operations work correctly
- [ ] All interactions work without full page reloads
- [ ] UI is responsive across mobile and desktop
- [ ] Error handling works correctly
- [ ] Bootstrap styling applied consistently
- [ ] Stimulus controllers work correctly

#### Go Best Practices
Not applicable (manual testing)

---

## Go Best Practices to Follow

### Code Organization
- **Reference**: `.claude/refs/go/best-practices.md` - Project Structure
- Use `internal/` for private packages
- Group related functionality into packages (model, store, service, handler, server)
- Keep packages small and focused on single responsibility

### Idiomatic Go
- **Reference**: `.claude/refs/go/idiomatic-go.md`
- Use constructor functions (New*) for initialization
- Accept interfaces, return structs
- Use named return values sparingly (only when it aids clarity)
- Keep functions short and focused

### Concurrency
- **Reference**: `.claude/refs/go/concurrency-patterns.md`
- Use `sync.RWMutex` for shared state protection
- Always defer unlock immediately after lock
- Prefer channels for communication, mutexes for state protection
- Return copies of slices to prevent concurrent modification

### Error Handling
- **Reference**: `.claude/refs/go/error-handling.md`
- Use sentinel errors for expected errors (ErrTaskNotFound, ErrEmptyTitle)
- Wrap errors with context using `fmt.Errorf` and `%w`
- Check errors with `errors.Is()` for wrapped errors
- Don't ignore errors (handle or log them)

### Design Patterns
- **Reference**: `.claude/refs/go/design-patterns.md`
- Use dependency injection via constructors
- Use interfaces for testing (could add TaskStore interface later)
- Keep handler, service, store layers separate

### Testing (Optional for Phase 1)
- **Reference**: `.claude/refs/go/testing-practices.md`
- Write table-driven tests
- Test concurrent access to TaskStore
- Use httptest for handler testing

## Validation Checklist

### Architecture Compliance
- [ ] All components from ARCHITECTURE.md implemented
- [ ] Directory structure matches architecture design
- [ ] API endpoints match specification
- [ ] Data model matches design (Task struct)
- [ ] Thread safety implemented correctly (RWMutex)

### Feature Requirements
- [ ] All functional requirements from FEATURE.md implemented
- [ ] All user stories acceptance criteria met
- [ ] All use cases working as described
- [ ] Non-functional requirements addressed (performance, security)

### Code Quality
- [ ] All files have package documentation
- [ ] Exported functions have godoc comments
- [ ] Error handling is comprehensive
- [ ] No obvious bugs or race conditions
- [ ] Code follows Go conventions

### User Experience
- [ ] Bootstrap styling applied correctly
- [ ] Responsive design works on mobile
- [ ] Error messages are clear and helpful
- [ ] Loading states provide feedback
- [ ] Completed tasks have visual distinction

## Dependencies Review

### Required Go Packages
- `github.com/gorilla/mux` - HTTP routing (already in go.mod)
- `html/template` - Template rendering (stdlib)
- `encoding/json` - JSON encoding/decoding (stdlib)
- `sync` - Mutex for thread safety (stdlib)
- `time` - Timestamps (stdlib)
- `strconv` - ID conversion (stdlib)
- `errors` - Error handling (stdlib)

### Frontend Dependencies
- **Bootstrap 5.3** - CSS framework (CDN)
- **Stimulus.js 3.2+** - Frontend framework (CDN via unpkg or esm.sh)

No additional Go dependencies needed beyond what's in the bootstrap template.

## Timeline Estimate

**Phase 1 (Foundation)**: 2-3 hours
- Clean up + directory structure + model + store

**Phase 2 (Backend Core)**: 3-4 hours
- Service + API handlers + Page handler + Router

**Phase 3 (Frontend)**: 2-3 hours
- Templates + Stimulus controller + Static files

**Phase 4 (Integration)**: 1-2 hours
- Wire up + Error handling + Config + Testing

**Total Estimated Time**: 8-12 hours of focused development

## Notes for Developer

1. **Start with Foundation**: Complete Phase 1 tasks in order (they build on each other)
2. **Backend and Frontend in Parallel**: Phase 2 and 3 can be done somewhat in parallel
3. **Test Early**: Don't wait until the end to test; test each component as you build it
4. **Reference Architecture**: Constantly refer back to ARCHITECTURE.md for implementation details
5. **Use Go Best Practices**: Follow the patterns from `.claude/refs/go/` documentation
6. **Keep It Simple**: Resist over-engineering; follow the architecture exactly
7. **Thread Safety**: Be extra careful with mutex usage (always defer unlock)
8. **Manual Testing**: Phase 4 includes comprehensive manual testing - don't skip it

## Success Criteria

This feature is complete when:
- ✅ All 15 tasks are checked off
- ✅ All acceptance criteria met
- ✅ Application runs with `make run`
- ✅ All CRUD operations work via browser
- ✅ UI is responsive and styled with Bootstrap
- ✅ No errors in browser console or server logs during normal operation
- ✅ Code follows Go best practices
- ✅ Manual testing completed successfully

---

<!-- COMPLETE -->
