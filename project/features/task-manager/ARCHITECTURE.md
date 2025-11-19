# Architecture: Simple Task Manager

## Overview

**Feature**: Simple Task Manager
**Architecture Type**: Monolithic Web Application
**Backend**: Go with Gorilla Mux
**Frontend**: Server-side rendering + Stimulus.js
**Storage**: In-memory (thread-safe)
**Created**: 2025-11-19
**Last Updated**: 2025-11-19

This document defines the technical architecture for a lightweight task management application that validates the agentic development workflow.

## Architecture Principles

### Core Principles

1. **Simplicity First**: Keep architecture minimal and focused
2. **Separation of Concerns**: Clear boundaries between layers
3. **Thread Safety**: Proper synchronization for concurrent access
4. **Progressive Enhancement**: Works without JavaScript, enhanced with it
5. **Idiomatic Go**: Follow Go best practices and conventions

### Technology Choices

| Component | Technology | Justification |
|-----------|-----------|---------------|
| HTTP Framework | Gorilla Mux | Already in bootstrap, familiar routing |
| Templating | Go html/template | Built-in, secure auto-escaping |
| Frontend Framework | Stimulus.js 3.2+ | Lightweight, progressive enhancement |
| CSS Framework | Bootstrap 5.3 | Responsive, component library |
| Storage | In-memory (sync.RWMutex) | Sufficient for testing, thread-safe |
| Data Format | JSON | Standard for API responses |

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│  ┌────────────────┐  ┌──────────────┐  ┌─────────────────┐ │
│  │  HTML/CSS      │  │ Stimulus.js  │  │   Bootstrap     │ │
│  │  (Templates)   │  │ Controllers  │  │   Styling       │ │
│  └────────────────┘  └──────────────┘  └─────────────────┘ │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/AJAX
┌──────────────────────────┴──────────────────────────────────┐
│                    Go HTTP Server                            │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  HTTP Layer                             │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐            │ │
│  │  │  Router  │  │ Handlers │  │Templates │            │ │
│  │  │  (Mux)   │  │          │  │          │            │ │
│  │  └──────────┘  └──────────┘  └──────────┘            │ │
│  └─────────────────────┬──────────────────────────────────┘ │
│  ┌─────────────────────┴──────────────────────────────────┐ │
│  │               Business Logic Layer                      │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │         TaskService (CRUD operations)            │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  └─────────────────────┬──────────────────────────────────┘ │
│  ┌─────────────────────┴──────────────────────────────────┐ │
│  │                Storage Layer                            │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │  TaskStore (In-Memory + Mutex)                   │ │ │
│  │  │  - tasks []Task                                  │ │ │
│  │  │  - mutex sync.RWMutex                            │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Request Flow

#### HTML Page Request
```
Browser GET /
    ↓
Router → PageHandler
    ↓
PageHandler.ServeTaskList()
    ├─ TaskService.GetAll()
    │  └─ TaskStore.GetAll() (read lock)
    ├─ Execute template with tasks
    └─ Return HTML response
```

#### AJAX API Request (Create Task)
```
Browser POST /api/tasks {"title": "..."}
    ↓
Router → APIHandler
    ↓
APIHandler.CreateTask()
    ├─ Validate input
    ├─ TaskService.Create(title)
    │  └─ TaskStore.Create(task) (write lock)
    ├─ Return JSON response
    └─ Status 201 Created
```

#### AJAX API Request (Toggle Task)
```
Browser PATCH /api/tasks/123/toggle
    ↓
Router → APIHandler
    ↓
APIHandler.ToggleTask(id)
    ├─ TaskService.Toggle(id)
    │  └─ TaskStore.Toggle(id) (write lock)
    ├─ Return JSON response
    └─ Status 200 OK
```

## Component Design

### 1. Data Layer

#### Task Model
```go
package model

import "time"

// Task represents a single task item
type Task struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"createdAt"`
}
```

#### Task Store (In-Memory)
```go
package store

import (
    "errors"
    "sync"
    "time"
)

// TaskStore provides thread-safe in-memory task storage
type TaskStore struct {
    tasks  []Task
    nextID int
    mu     sync.RWMutex
}

func NewTaskStore() *TaskStore {
    return &TaskStore{
        tasks:  make([]Task, 0),
        nextID: 1,
    }
}

// GetAll returns all tasks (uses read lock)
func (s *TaskStore) GetAll() []Task

// GetByID returns a single task (uses read lock)
func (s *TaskStore) GetByID(id string) (Task, error)

// Create adds a new task (uses write lock)
func (s *TaskStore) Create(title string) Task

// Toggle changes completion status (uses write lock)
func (s *TaskStore) Toggle(id string) (Task, error)

// Delete removes a task (uses write lock)
func (s *TaskStore) Delete(id string) error
```

**Thread Safety Strategy:**
- Use `sync.RWMutex` for concurrent access
- Read operations use `RLock()` (multiple readers allowed)
- Write operations use `Lock()` (exclusive access)
- Always defer unlock to prevent deadlocks

### 2. Business Logic Layer

#### Task Service
```go
package service

// TaskService handles business logic for tasks
type TaskService struct {
    store *store.TaskStore
}

func NewTaskService(store *store.TaskStore) *TaskService {
    return &TaskService{store: store}
}

// GetAll retrieves all tasks
func (s *TaskService) GetAll() []Task

// Create creates a new task with validation
func (s *TaskService) Create(title string) (Task, error) {
    // Validate title (not empty, max length)
    // Call store.Create
}

// Toggle toggles task completion status
func (s *TaskService) Toggle(id string) (Task, error)

// Delete removes a task
func (s *TaskService) Delete(id string) error
```

**Responsibilities:**
- Input validation (title not empty, length limits)
- Business rules enforcement
- Error handling and messaging
- Logging

### 3. HTTP Layer

#### Router Configuration
```go
package server

import "github.com/gorilla/mux"

func NewRouter(pageHandler *PageHandler, apiHandler *APIHandler) *mux.Router {
    r := mux.NewRouter()

    // Static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
        http.FileServer(http.Dir("static"))))

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

#### Page Handler
```go
package handler

import (
    "html/template"
    "net/http"
)

// PageHandler handles HTML page requests
type PageHandler struct {
    service   *service.TaskService
    templates *template.Template
}

func NewPageHandler(service *service.TaskService) *PageHandler {
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
        Tasks []Task
    }{
        Tasks: tasks,
    }

    if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
```

#### API Handler
```go
package handler

import (
    "encoding/json"
    "net/http"
)

// APIHandler handles JSON API requests
type APIHandler struct {
    service *service.TaskService
}

func NewAPIHandler(service *service.TaskService) *APIHandler {
    return &APIHandler{service: service}
}

// GetTasks returns all tasks as JSON
func (h *APIHandler) GetTasks(w http.ResponseWriter, r *http.Request)

// CreateTask creates a new task from JSON
func (h *APIHandler) CreateTask(w http.ResponseWriter, r *http.Request)

// ToggleTask toggles task completion status
func (h *APIHandler) ToggleTask(w http.ResponseWriter, r *http.Request)

// DeleteTask deletes a task
func (h *APIHandler) DeleteTask(w http.ResponseWriter, r *http.Request)
```

**Error Responses:**
```json
{
  "error": "Task not found",
  "code": "NOT_FOUND"
}
```

**Success Responses:**
```json
{
  "id": "1",
  "title": "Buy groceries",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z"
}
```

### 4. Template Layer

#### Directory Structure
```
templates/
├── base.html          # Base layout with <!DOCTYPE>, <head>, <body>
├── index.html         # Main task list page
└── partials/
    ├── task-item.html # Individual task row
    └── task-form.html # Task creation form
```

#### Base Template
```html
<!DOCTYPE html>
<html lang="en" data-controller="tasks">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Task Manager</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        {{template "content" .}}
    </div>

    <script type="module" src="/static/js/app.js"></script>
</body>
</html>
```

#### Index Template
```html
{{define "content"}}
<div class="row">
    <div class="col-md-8 offset-md-2">
        <h1 class="mb-4">Task Manager</h1>

        <!-- Task Form -->
        {{template "task-form" .}}

        <!-- Task List -->
        <div data-tasks-target="list">
            {{range .Tasks}}
                {{template "task-item" .}}
            {{else}}
                <p class="text-muted">No tasks yet. Add one above!</p>
            {{end}}
        </div>
    </div>
</div>
{{end}}
```

### 5. Frontend Layer (Stimulus.js)

#### Controller Structure
```javascript
// static/js/controllers/tasks_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = ["list", "form", "input"]

  // Add new task
  async create(event) {
    event.preventDefault()
    const title = this.inputTarget.value.trim()

    if (!title) return

    const response = await fetch('/api/tasks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title })
    })

    if (response.ok) {
      const task = await response.json()
      this.addTaskToList(task)
      this.inputTarget.value = ''
    }
  }

  // Toggle task completion
  async toggle(event) {
    const taskId = event.target.dataset.taskId

    const response = await fetch(`/api/tasks/${taskId}/toggle`, {
      method: 'PATCH'
    })

    if (response.ok) {
      const task = await response.json()
      this.updateTaskInList(task)
    }
  }

  // Delete task
  async delete(event) {
    if (!confirm('Delete this task?')) return

    const taskId = event.target.dataset.taskId

    const response = await fetch(`/api/tasks/${taskId}`, {
      method: 'DELETE'
    })

    if (response.ok) {
      this.removeTaskFromList(taskId)
    }
  }

  // Helper methods for DOM manipulation
  addTaskToList(task) { ... }
  updateTaskInList(task) { ... }
  removeTaskFromList(taskId) { ... }
}
```

#### Application Setup
```javascript
// static/js/app.js
import { Application } from "@hotwired/stimulus"
import TasksController from "./controllers/tasks_controller.js"

const application = Application.start()
application.register("tasks", TasksController)
```

## Directory Structure

```
apps/test-task-manager/
├── cmd/
│   └── test-task-manager/
│       └── main.go                 # Application entry point
├── internal/
│   ├── model/
│   │   └── task.go                 # Task model
│   ├── store/
│   │   └── task_store.go           # In-memory storage
│   ├── service/
│   │   └── task_service.go         # Business logic
│   ├── handler/
│   │   ├── page_handler.go         # HTML handlers
│   │   └── api_handler.go          # API handlers
│   └── server/
│       ├── server.go                # HTTP server setup
│       └── routes.go                # Route configuration
├── templates/
│   ├── base.html                    # Base layout
│   ├── index.html                   # Main page
│   └── partials/
│       ├── task-item.html           # Task item partial
│       └── task-form.html           # Task form partial
├── static/
│   ├── js/
│   │   ├── app.js                   # Stimulus app setup
│   │   └── controllers/
│   │       └── tasks_controller.js  # Tasks controller
│   └── css/
│       └── custom.css               # Custom styles (if needed)
├── go.mod                            # Go module definition
├── go.sum                            # Go dependencies
├── Makefile                          # Build automation
└── README.md                         # Project documentation
```

## API Specification

### Endpoints

#### GET /
**Description**: Render main task list page
**Response**: HTML document
**Status Codes**:
- 200 OK: Page rendered successfully
- 500 Internal Server Error: Template rendering failed

#### GET /api/tasks
**Description**: Get all tasks as JSON
**Response**:
```json
[
  {
    "id": "1",
    "title": "Buy groceries",
    "completed": false,
    "createdAt": "2025-11-19T10:00:00Z"
  },
  {
    "id": "2",
    "title": "Walk the dog",
    "completed": true,
    "createdAt": "2025-11-19T09:00:00Z"
  }
]
```
**Status Codes**:
- 200 OK: Tasks retrieved successfully

#### POST /api/tasks
**Description**: Create a new task
**Request**:
```json
{
  "title": "Buy groceries"
}
```
**Response**:
```json
{
  "id": "1",
  "title": "Buy groceries",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z"
}
```
**Status Codes**:
- 201 Created: Task created successfully
- 400 Bad Request: Invalid input (empty title, too long, etc.)
- 500 Internal Server Error: Creation failed

**Validation Rules**:
- Title must not be empty
- Title must not exceed 255 characters
- Title must be trimmed of whitespace

#### PATCH /api/tasks/{id}/toggle
**Description**: Toggle task completion status
**URL Parameters**:
- `id`: Task ID
**Response**:
```json
{
  "id": "1",
  "title": "Buy groceries",
  "completed": true,
  "createdAt": "2025-11-19T10:00:00Z"
}
```
**Status Codes**:
- 200 OK: Task toggled successfully
- 404 Not Found: Task ID doesn't exist
- 500 Internal Server Error: Toggle failed

#### DELETE /api/tasks/{id}
**Description**: Delete a task
**URL Parameters**:
- `id`: Task ID
**Response**:
```json
{
  "message": "Task deleted successfully"
}
```
**Status Codes**:
- 200 OK: Task deleted successfully
- 404 Not Found: Task ID doesn't exist
- 500 Internal Server Error: Deletion failed

## Data Flow Patterns

### Create Task Flow
```
1. User types title in form, presses Enter
   ↓
2. Stimulus controller intercepts submit event
   ↓
3. Controller validates input client-side
   ↓
4. Controller sends POST /api/tasks with JSON
   ↓
5. APIHandler.CreateTask() receives request
   ↓
6. Handler decodes JSON, validates input
   ↓
7. Handler calls TaskService.Create(title)
   ↓
8. Service validates business rules
   ↓
9. Service calls TaskStore.Create(title)
   ↓
10. Store acquires write lock
   ↓
11. Store creates task with new ID
   ↓
12. Store adds task to slice
   ↓
13. Store releases write lock
   ↓
14. Store returns created task
   ↓
15. Service returns task to handler
   ↓
16. Handler encodes task as JSON
   ↓
17. Handler sends 201 Created response
   ↓
18. Controller receives response
   ↓
19. Controller adds task to DOM
   ↓
20. Controller clears form input
   ↓
21. User sees new task in list
```

### Toggle Task Flow
```
1. User clicks checkbox/button
   ↓
2. Stimulus controller captures click event
   ↓
3. Controller sends PATCH /api/tasks/{id}/toggle
   ↓
4. APIHandler.ToggleTask() receives request
   ↓
5. Handler extracts task ID from URL
   ↓
6. Handler calls TaskService.Toggle(id)
   ↓
7. Service calls TaskStore.Toggle(id)
   ↓
8. Store acquires write lock
   ↓
9. Store finds task by ID
   ↓
10. Store toggles completed field
   ↓
11. Store releases write lock
   ↓
12. Store returns updated task
   ↓
13. Service returns task to handler
   ↓
14. Handler encodes task as JSON
   ↓
15. Handler sends 200 OK response
   ↓
16. Controller receives response
   ↓
17. Controller updates task styling in DOM
   ↓
18. User sees updated completion status
```

## Security Considerations

### Input Validation
- **Server-side validation required** for all inputs
- Title sanitization to prevent XSS
- Maximum length enforcement (255 chars)
- HTML escaping via Go templates

### Thread Safety
- **RWMutex** protects shared state
- Read operations use RLock (concurrent reads OK)
- Write operations use Lock (exclusive access)
- Always defer unlock to prevent deadlocks

### HTTP Security
- Content-Type headers set correctly
- CORS not needed (same-origin)
- No CSRF protection needed for testing (can add later)
- Input sanitization prevents injection attacks

## Performance Considerations

### In-Memory Storage
- **O(1)** time complexity for create operations
- **O(n)** time complexity for search/delete (linear scan)
- **Acceptable** for testing with < 1000 tasks
- **Trade-off**: Speed vs. persistence

### Concurrent Access
- Multiple readers can read simultaneously (RLock)
- Single writer blocks all access (Lock)
- Expected load: ~50 concurrent users
- Bottleneck: Write operations (create, toggle, delete)

### Frontend Performance
- Minimal JavaScript payload (Stimulus.js ~33KB gzipped)
- Bootstrap loaded from CDN (cached)
- No unnecessary re-renders (targeted DOM updates)
- AJAX requests prevent full page reloads

## Error Handling Strategy

### HTTP Errors
```go
// Consistent error response format
type ErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}

// Error codes
const (
    ErrNotFound       = "NOT_FOUND"
    ErrInvalidInput   = "INVALID_INPUT"
    ErrInternalServer = "INTERNAL_SERVER_ERROR"
)
```

### Logging Strategy
- Log all HTTP requests (method, path, status, duration)
- Log validation failures at INFO level
- Log errors at ERROR level with stack trace
- Log task operations at DEBUG level

### Graceful Degradation
- If JavaScript fails, page still displays tasks (server-rendered)
- Form submission can work without JavaScript (add non-AJAX fallback)
- Bootstrap ensures mobile responsiveness

## Deployment Configuration

### Environment Variables
```bash
# Server configuration
HTTP_PORT=8080              # Port for HTTP server
LOG_LEVEL=debug             # Logging verbosity

# Not needed for task manager (remove from .env)
# DATABASE_URL                # No database
# SENTRY_DSN                  # No error tracking
# PUBSUB_EMULATOR             # No messaging
```

### Build Configuration
```makefile
# Makefile updates
CMD=go run ./cmd/test-task-manager/main.go -port=8080

run:
	${CMD}

build:
	go build -o bin/test-task-manager ./cmd/test-task-manager

test:
	go test ./internal/...
```

### Server Initialization
```go
// cmd/test-task-manager/main.go
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

    // Start server
    port := os.Getenv("HTTP_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}
```

## Testing Strategy

### Unit Tests (Developer will implement)
```go
// internal/store/task_store_test.go
func TestTaskStore_Create(t *testing.T)
func TestTaskStore_GetAll(t *testing.T)
func TestTaskStore_Toggle(t *testing.T)
func TestTaskStore_Delete(t *testing.T)
func TestTaskStore_ConcurrentAccess(t *testing.T)  // Test thread safety

// internal/service/task_service_test.go
func TestTaskService_Validation(t *testing.T)
func TestTaskService_Create(t *testing.T)

// internal/handler/api_handler_test.go
func TestAPIHandler_CreateTask(t *testing.T)
func TestAPIHandler_ToggleTask(t *testing.T)
```

### Integration Tests
```go
// Test full request/response cycle
func TestCreateTaskEndToEnd(t *testing.T)
func TestToggleTaskEndToEnd(t *testing.T)
func TestDeleteTaskEndToEnd(t *testing.T)
```

### Manual Testing
- User will perform manual testing in browser
- Test all CRUD operations
- Test concurrent access (multiple browser tabs)
- Test error scenarios (empty title, invalid ID)

## Optional Enhancements (Phase 2)

### Task Filtering
**Client-side implementation:**
```javascript
// Add to tasks_controller.js
filter(event) {
  const filter = event.target.dataset.filter  // 'all', 'active', 'completed'
  const tasks = this.listTarget.querySelectorAll('[data-task-id]')

  tasks.forEach(task => {
    const isCompleted = task.classList.contains('completed')

    if (filter === 'all') {
      task.style.display = ''
    } else if (filter === 'active' && !isCompleted) {
      task.style.display = ''
    } else if (filter === 'completed' && isCompleted) {
      task.style.display = ''
    } else {
      task.style.display = 'none'
    }
  })
}
```

### Task Counters
```javascript
// Add computed properties
get totalCount() {
  return this.listTarget.querySelectorAll('[data-task-id]').length
}

get activeCount() {
  return this.listTarget.querySelectorAll('[data-task-id]:not(.completed)').length
}

get completedCount() {
  return this.listTarget.querySelectorAll('[data-task-id].completed').length
}
```

## Implementation Notes

### Simplifications from Bootstrap Template
- **Remove database**: No MySQL, Cloud SQL, or migrations
- **Remove Pub/Sub**: No messaging infrastructure
- **Remove Sentry**: No error tracking (use console logging)
- **Remove app lifecycle**: Simplified initialization
- **Keep Gorilla Mux**: For routing
- **Keep logging**: Basic request/error logging

### Key Implementation Tasks
1. Create task model and store
2. Implement thread-safe storage with RWMutex
3. Build HTTP handlers (page + API)
4. Create HTML templates with Bootstrap
5. Implement Stimulus.js controller
6. Set up routing and static file serving
7. Add input validation
8. Implement error handling
9. Add logging
10. Manual testing

## Architecture Review Checklist

- [x] All functional requirements addressed
- [x] Thread safety ensured (RWMutex)
- [x] Clear separation of concerns (layers)
- [x] RESTful API design
- [x] Progressive enhancement (works without JS)
- [x] Bootstrap integration defined
- [x] Stimulus.js integration defined
- [x] Error handling strategy defined
- [x] Security considerations addressed
- [x] Performance considerations addressed
- [x] Directory structure defined
- [x] API specification complete
- [x] Deployment configuration defined

## Open Questions

None. Architecture is clear and ready for implementation planning.

## Assumptions

- Single-instance deployment (no load balancing needed)
- In-memory storage acceptable (no persistence required)
- Modern browser support (ES6+ JavaScript)
- No authentication/authorization needed
- Testing environment only (not production-ready)

## Related Documentation

- **FEATURE.md**: Business requirements and user stories
- **TODO.md**: Will contain detailed implementation tasks
- **Go Best Practices**: `.claude/refs/go/best-practices.md`
- **Idiomatic Go**: `.claude/refs/go/idiomatic-go.md`

---

<!-- COMPLETE -->
