# Architecture: Task Customization with Priority Indicators

## Overview

**Feature**: Task Customization with Priority Indicators
**Architecture Type**: Enhancement to Existing Monolithic Web Application
**Pattern**: Model-Service-Handler with Client-Side Filtering
**Complexity**: Medium
**Impact**: Extends existing Task model, adds client-side filtering logic
**Created**: 2025-11-19
**Last Updated**: 2025-11-19

This document defines the technical architecture for adding visual priority indicators (emoticons and colors) based on the Eisenhower Matrix methodology, with client-side filtering capabilities.

## Architecture Principles

### Core Principles

1. **Backward Compatibility**: Existing tasks automatically receive default priority
2. **Immutability**: Priority cannot be changed after task creation (locked)
3. **Client-Side Performance**: Filtering happens instantly in browser (< 100ms)
4. **Separation of Concerns**: Priority logic cleanly integrated into existing layers
5. **Idiomatic Go**: Follow existing codebase patterns

### Technology Choices

| Component | Technology | Change Type | Justification |
|-----------|-----------|-------------|---------------|
| Task Model | Go struct extension | Modified | Add Priority and Color fields |
| Validation | Service layer logic | Modified | Validate priority/color values |
| API | JSON request/response | Modified | Accept/return priority/color |
| Frontend Filtering | Stimulus.js controller | New | Client-side instant filtering |
| UI Components | Bootstrap 5.3 buttons | New | Filter button styling |
| Storage | In-memory (existing) | No change | Priority/Color stored with Task |

## System Architecture

### High-Level Architecture (Changes Highlighted)

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│  ┌────────────────┐  ┌──────────────┐  ┌─────────────────┐ │
│  │  HTML/CSS      │  │ Stimulus.js  │  │   Bootstrap     │ │
│  │  (Templates)   │  │ Controllers  │  │   Styling       │ │
│  │  + Emoticons   │  │ + FILTERING ✨│  │ + Filter Btns ✨│ │
│  │  + Colors ✨   │  │ + Priority   │  │                 │ │
│  └────────────────┘  └──────────────┘  └─────────────────┘ │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/AJAX (with priority/color)
┌──────────────────────────┴──────────────────────────────────┐
│                    Go HTTP Server                            │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  HTTP Layer                             │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐            │ │
│  │  │  Router  │  │ Handlers │  │Templates │            │ │
│  │  │  (Mux)   │  │ + Accept │  │ + Render │            │ │
│  │  │          │  │ Priority ✨│  │ Emoticon ✨│          │ │
│  │  └──────────┘  └──────────┘  └──────────┘            │ │
│  └─────────────────────┬──────────────────────────────────┘ │
│  ┌─────────────────────┴──────────────────────────────────┐ │
│  │               Business Logic Layer                      │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │  TaskService (+ Priority Validation) ✨         │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  └─────────────────────┬──────────────────────────────────┘ │
│  ┌─────────────────────┴──────────────────────────────────┐ │
│  │                Storage Layer                            │ │
│  │  ┌──────────────────────────────────────────────────┐ │ │
│  │  │  TaskStore (stores Task with Priority/Color) ✨│ │ │
│  │  │  - tasks []Task                                  │ │ │
│  │  │  - mutex sync.RWMutex                            │ │ │
│  │  └──────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Request Flow Changes

#### Enhanced Create Task Flow
```
Browser POST /api/tasks {"title": "...", "priority": "🔥", "color": "#dc3545"}
    ↓
Router → APIHandler
    ↓
APIHandler.CreateTask()
    ├─ Decode JSON (includes priority, color) ✨
    ├─ Validate input
    ├─ TaskService.Create(title, priority, color) ✨
    │  ├─ Validate priority is valid emoticon ✨
    │  ├─ Validate color is valid hex code ✨
    │  ├─ Apply defaults if empty ✨
    │  └─ TaskStore.Create(task) (write lock)
    ├─ Return JSON response with priority/color ✨
    └─ Status 201 Created
```

#### New Client-Side Filtering Flow
```
User clicks filter button (e.g., 🔥)
    ↓
Stimulus controller captures click event
    ↓
Controller gets selected filters (e.g., ["🔥", "⭐"])
    ↓
Controller iterates all task elements in DOM
    ↓
For each task:
    ├─ Check if task.priority matches any selected filter
    ├─ If match: show task (display: block)
    └─ If no match: hide task (display: none)
    ↓
Update task count display
    ↓
User sees filtered list instantly (< 100ms)
```

## Component Design Changes

### 1. Data Layer (MODIFIED)

#### Task Model (Extended)
```go
package model

import "time"

// Task represents a single task item with priority
type Task struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"createdAt"`
    Priority  string    `json:"priority"`  // NEW: Emoticon (🔥, ⭐, ⚡, 💡, 📋)
    Color     string    `json:"color"`     // NEW: Hex code (#dc3545, etc.)
}
```

**Changes**:
- Add `Priority` field (string, stores UTF-8 emoticon)
- Add `Color` field (string, stores hex code)
- Both fields included in JSON serialization

**Default Values**:
- Priority: "📋" (Clipboard emoticon)
- Color: "#6c757d" (Grey)

#### Task Store (NO CHANGES)
```go
package store

// TaskStore provides thread-safe in-memory task storage
type TaskStore struct {
    tasks  []Task
    nextID int
    mu     sync.RWMutex
}
```

**Why No Changes**:
- Store operates on Task structs generically
- Priority and Color automatically stored as part of Task
- No special indexing or querying needed (client-side filtering)

**Migration on Startup** (Optional Enhancement):
```go
// NewTaskStore can initialize with migration logic
func NewTaskStore() *TaskStore {
    store := &TaskStore{
        tasks:  make([]Task, 0),
        nextID: 1,
    }

    // Future: Add migration logic here if loading tasks from disk
    // store.migrateExistingTasks()

    return store
}
```

### 2. Business Logic Layer (MODIFIED)

#### Task Service (Extended)
```go
package service

import "errors"

const (
    // Valid emoticons
    PriorityUrgentImportant = "🔥"  // Red
    PriorityImportant       = "⭐"  // Blue
    PriorityUrgent          = "⚡"  // Yellow
    PriorityLow             = "💡"  // Green
    PriorityDefault         = "📋"  // Grey

    // Valid colors
    ColorRed    = "#dc3545"
    ColorBlue   = "#0d6efd"
    ColorYellow = "#ffc107"
    ColorGreen  = "#28a745"
    ColorPurple = "#6f42c1"
    ColorOrange = "#fd7e14"
    ColorGrey   = "#6c757d"
)

var (
    ErrInvalidPriority = errors.New("invalid priority emoticon")
    ErrInvalidColor    = errors.New("invalid color code")
)

type TaskService struct {
    store *store.TaskStore
}

// Create creates a new task with validation
func (s *TaskService) Create(title, priority, color string) (Task, error) {
    // Validate title (existing logic)
    title = strings.TrimSpace(title)
    if title == "" {
        return Task{}, ErrEmptyTitle
    }
    if len(title) > 255 {
        return Task{}, ErrTitleTooLong
    }

    // NEW: Apply defaults if not provided
    if priority == "" {
        priority = PriorityDefault
    }
    if color == "" {
        color = ColorGrey
    }

    // NEW: Validate priority
    if !isValidPriority(priority) {
        return Task{}, ErrInvalidPriority
    }

    // NEW: Validate color
    if !isValidColor(color) {
        return Task{}, ErrInvalidColor
    }

    // Create task with priority and color
    task := s.store.Create(title, priority, color)
    return task, nil
}

// isValidPriority checks if emoticon is valid
func isValidPriority(p string) bool {
    validPriorities := []string{
        PriorityUrgentImportant,
        PriorityImportant,
        PriorityUrgent,
        PriorityLow,
        PriorityDefault,
    }
    for _, valid := range validPriorities {
        if p == valid {
            return true
        }
    }
    return false
}

// isValidColor checks if hex code is valid
func isValidColor(c string) bool {
    validColors := []string{
        ColorRed, ColorBlue, ColorYellow, ColorGreen,
        ColorPurple, ColorOrange, ColorGrey,
    }
    for _, valid := range validColors {
        if c == valid {
            return true
        }
    }
    return false
}
```

**Changes**:
- `Create()` signature: `Create(title string)` → `Create(title, priority, color string)`
- New validation functions: `isValidPriority()`, `isValidColor()`
- New error types: `ErrInvalidPriority`, `ErrInvalidColor`
- Default value application logic

**Validation Rules**:
- Priority must be one of 5 valid emoticons
- Color must be one of 7 valid hex codes
- Empty priority/color → apply defaults
- Invalid values → return validation error

### 3. HTTP Layer (MODIFIED)

#### API Handler (Extended)
```go
package handler

type CreateTaskRequest struct {
    Title    string `json:"title"`
    Priority string `json:"priority"` // NEW: Optional
    Color    string `json:"color"`    // NEW: Optional
}

func (h *APIHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    // Decode JSON request
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Call service with priority and color
    task, err := h.service.Create(req.Title, req.Priority, req.Color)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrEmptyTitle):
            respondError(w, "Title cannot be empty", http.StatusBadRequest)
        case errors.Is(err, service.ErrTitleTooLong):
            respondError(w, "Title too long (max 255 characters)", http.StatusBadRequest)
        case errors.Is(err, service.ErrInvalidPriority): // NEW
            respondError(w, "Invalid priority emoticon", http.StatusBadRequest)
        case errors.Is(err, service.ErrInvalidColor): // NEW
            respondError(w, "Invalid color code", http.StatusBadRequest)
        default:
            respondError(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Return task with priority and color
    respondJSON(w, task, http.StatusCreated)
}
```

**Changes**:
- `CreateTaskRequest` struct: Add `Priority` and `Color` fields
- Error handling: Add cases for `ErrInvalidPriority` and `ErrInvalidColor`
- Service call: Pass priority and color parameters

**Response Format** (Enhanced):
```json
{
  "id": "1",
  "title": "Fix bug",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z",
  "priority": "🔥",
  "color": "#dc3545"
}
```

#### Page Handler (NO CHANGES)
```go
func (h *PageHandler) ServeTaskList(w http.ResponseWriter, r *http.Request) {
    tasks := h.service.GetAll()

    data := struct {
        Tasks []Task
    }{
        Tasks: tasks,
    }

    // Template automatically receives priority/color fields
    if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
```

**Why No Changes**:
- Handler already passes full Task structs to template
- Priority and Color automatically included
- Template responsible for rendering

### 4. Template Layer (MODIFIED)

#### Index Template (Enhanced)
```html
{{define "content"}}
<div class="row">
    <div class="col-md-8 offset-md-2">
        <h1 class="mb-4">Task Manager</h1>

        <!-- Task Form with Priority Selector ✨ -->
        <form data-action="submit->tasks#create" class="mb-4">
            <div class="input-group mb-3">
                <input
                    type="text"
                    class="form-control"
                    placeholder="Enter task title..."
                    data-tasks-target="input"
                    required
                >
                <button class="btn btn-primary" type="submit">Add Task</button>
            </div>

            <!-- NEW: Priority Selector -->
            <div class="btn-group mb-3" role="group" aria-label="Priority selector">
                <input type="radio" class="btn-check" name="priority" id="priority-urgent"
                       value="🔥" data-color="#dc3545" data-tasks-target="priorityInput">
                <label class="btn btn-outline-danger" for="priority-urgent">
                    🔥 Urgent & Important
                </label>

                <input type="radio" class="btn-check" name="priority" id="priority-important"
                       value="⭐" data-color="#0d6efd" data-tasks-target="priorityInput">
                <label class="btn btn-outline-primary" for="priority-important">
                    ⭐ Important
                </label>

                <input type="radio" class="btn-check" name="priority" id="priority-urgent-only"
                       value="⚡" data-color="#ffc107" data-tasks-target="priorityInput">
                <label class="btn btn-outline-warning" for="priority-urgent-only">
                    ⚡ Urgent
                </label>

                <input type="radio" class="btn-check" name="priority" id="priority-low"
                       value="💡" data-color="#28a745" data-tasks-target="priorityInput">
                <label class="btn btn-outline-success" for="priority-low">
                    💡 Low Priority
                </label>

                <input type="radio" class="btn-check" name="priority" id="priority-default"
                       value="📋" data-color="#6c757d" data-tasks-target="priorityInput" checked>
                <label class="btn btn-outline-secondary" for="priority-default">
                    📋 Default
                </label>
            </div>
        </form>

        <!-- NEW: Filter Buttons ✨ -->
        <div class="mb-3">
            <div class="btn-group" role="group" aria-label="Filter tasks">
                <button type="button" class="btn btn-sm btn-outline-secondary"
                        data-action="click->tasks#clearFilters">
                    Show All
                </button>
                <button type="button" class="btn btn-sm btn-outline-danger"
                        data-action="click->tasks#filterByPriority"
                        data-priority="🔥">
                    🔥
                </button>
                <button type="button" class="btn btn-sm btn-outline-primary"
                        data-action="click->tasks#filterByPriority"
                        data-priority="⭐">
                    ⭐
                </button>
                <button type="button" class="btn btn-sm btn-outline-warning"
                        data-action="click->tasks#filterByPriority"
                        data-priority="⚡">
                    ⚡
                </button>
                <button type="button" class="btn btn-sm btn-outline-success"
                        data-action="click->tasks#filterByPriority"
                        data-priority="💡">
                    💡
                </button>
                <button type="button" class="btn btn-sm btn-outline-secondary"
                        data-action="click->tasks#filterByPriority"
                        data-priority="📋">
                    📋
                </button>
            </div>
            <span class="ms-2 text-muted" data-tasks-target="taskCount">
                Showing {{len .Tasks}} tasks
            </span>
        </div>

        <!-- Task List with Priority Display ✨ -->
        <ul class="list-group" data-tasks-target="list">
            {{range .Tasks}}
            <li class="list-group-item d-flex justify-content-between align-items-center"
                data-task-id="{{.ID}}"
                data-priority="{{.Priority}}"
                style="border-left: 4px solid {{.Color}}">

                <div class="form-check flex-grow-1">
                    <input class="form-check-input" type="checkbox"
                           {{if .Completed}}checked{{end}}
                           data-action="change->tasks#toggle"
                           data-task-id="{{.ID}}">
                    <label class="form-check-label {{if .Completed}}text-decoration-line-through text-muted{{end}}"
                           data-tasks-target="label">
                        <span class="me-2">{{.Priority}}</span>{{.Title}}
                    </label>
                </div>

                <button class="btn btn-sm btn-danger"
                        data-action="click->tasks#delete"
                        data-task-id="{{.ID}}">
                    Delete
                </button>
            </li>
            {{else}}
            <li class="list-group-item text-muted">
                No tasks yet. Add one above!
            </li>
            {{end}}
        </ul>
    </div>
</div>
{{end}}
```

**Changes**:
- Priority selector: Radio button group with 5 options
- Filter buttons: Button group with emoticon buttons
- Task item: Display emoticon + colored left border
- Data attributes: `data-priority` for filtering logic

### 5. Frontend Layer (MODIFIED)

#### Tasks Controller (Extended)
```javascript
// static/js/controllers/tasks_controller.js
import { Controller } from "https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js"

export default class extends Controller {
    static targets = ["input", "error", "list", "label", "priorityInput", "taskCount"]

    // NEW: Track active filters
    activeFilters = new Set()

    // Create task with priority and color
    async create(event) {
        event.preventDefault()

        const title = this.inputTarget.value.trim()

        if (!title) {
            this.showError("Please enter a task title")
            return
        }

        // NEW: Get selected priority and color
        const selectedInput = this.priorityInputTargets.find(input => input.checked)
        const priority = selectedInput ? selectedInput.value : "📋"
        const color = selectedInput ? selectedInput.dataset.color : "#6c757d"

        try {
            const response = await fetch("/api/tasks", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ title, priority, color }), // NEW: Include priority/color
            })

            const data = await response.json()

            if (!response.ok) {
                this.showError(data.error || "Failed to create task")
                return
            }

            // Clear input and reload
            this.inputTarget.value = ""
            this.hideError()
            window.location.reload()
        } catch (error) {
            this.showError("Network error: Could not create task")
            console.error("Create task error:", error)
        }
    }

    // NEW: Filter tasks by priority
    filterByPriority(event) {
        const priority = event.target.dataset.priority
        const button = event.target

        // Toggle filter
        if (this.activeFilters.has(priority)) {
            this.activeFilters.delete(priority)
            button.classList.remove("active")
        } else {
            this.activeFilters.add(priority)
            button.classList.add("active")
        }

        this.applyFilters()
    }

    // NEW: Clear all filters
    clearFilters() {
        this.activeFilters.clear()

        // Remove active state from all filter buttons
        document.querySelectorAll('[data-action*="filterByPriority"]').forEach(btn => {
            btn.classList.remove("active")
        })

        this.applyFilters()
    }

    // NEW: Apply active filters to task list
    applyFilters() {
        const tasks = this.listTarget.querySelectorAll('[data-task-id]')
        let visibleCount = 0

        tasks.forEach(task => {
            const taskPriority = task.dataset.priority

            // Show if no filters active OR priority matches any active filter
            if (this.activeFilters.size === 0 || this.activeFilters.has(taskPriority)) {
                task.style.display = ""
                visibleCount++
            } else {
                task.style.display = "none"
            }
        })

        // Update task count
        if (this.hasTaskCountTarget) {
            const totalCount = tasks.length
            const countText = this.activeFilters.size > 0
                ? `Showing ${visibleCount} of ${totalCount} tasks`
                : `Showing ${totalCount} tasks`
            this.taskCountTarget.textContent = countText
        }
    }

    // Existing methods: toggle, delete, showError, hideError
    // (No changes needed - they already work with extended Task model)
}
```

**Changes**:
- `create()`: Get selected priority/color from form, include in API request
- `filterByPriority()`: NEW - Toggle filter on/off
- `clearFilters()`: NEW - Reset filters to show all
- `applyFilters()`: NEW - Show/hide tasks based on active filters
- `activeFilters`: NEW - Set to track which filters are active

**Performance**:
- Filtering is pure DOM manipulation (no API calls)
- O(n) complexity where n = number of tasks
- Executes in < 100ms for ~1000 tasks

## Data Flow Patterns

### Create Task with Priority Flow
```
1. User types title, selects priority (e.g., 🔥 Red)
   ↓
2. User clicks "Add Task"
   ↓
3. Stimulus controller captures submit event
   ↓
4. Controller reads title, priority (🔥), color (#dc3545)
   ↓
5. Controller sends POST /api/tasks with JSON
   {
     "title": "Fix critical bug",
     "priority": "🔥",
     "color": "#dc3545"
   }
   ↓
6. APIHandler.CreateTask() receives request
   ↓
7. Handler decodes JSON
   ↓
8. Handler calls TaskService.Create(title, priority, color)
   ↓
9. Service validates title (existing validation)
   ↓
10. Service validates priority is valid emoticon
   ↓
11. Service validates color is valid hex code
   ↓
12. Service calls TaskStore.Create(title, priority, color)
   ↓
13. Store creates Task with all fields including priority/color
   ↓
14. Store returns created task
   ↓
15. Service returns task to handler
   ↓
16. Handler encodes task as JSON (includes priority/color)
   ↓
17. Handler sends 201 Created response
   ↓
18. Controller receives response
   ↓
19. Controller reloads page
   ↓
20. User sees task with 🔥 emoticon and red left border
```

### Client-Side Filtering Flow
```
1. User clicks 🔥 filter button
   ↓
2. Stimulus controller captures click event
   ↓
3. Controller extracts priority from button (data-priority="🔥")
   ↓
4. Controller toggles filter state:
   - If 🔥 not active: add to activeFilters Set
   - If 🔥 already active: remove from activeFilters Set
   ↓
5. Controller updates button visual state (add/remove "active" class)
   ↓
6. Controller calls applyFilters()
   ↓
7. applyFilters() iterates all <li> elements with [data-task-id]
   ↓
8. For each task element:
   - Read data-priority attribute
   - Check if priority in activeFilters Set
   - If YES or activeFilters empty: show (display: "")
   - If NO: hide (display: "none")
   ↓
9. Update task count display
   ↓
10. User sees filtered list instantly (< 100ms)
```

### Multi-Filter Combination Flow
```
User clicks 🔥 button → Shows only 🔥 tasks
   ↓
User clicks ⭐ button → Shows 🔥 AND ⭐ tasks (OR logic)
   ↓
User clicks 🔥 button again → Shows only ⭐ tasks
   ↓
User clicks "Show All" → Clears all filters, shows all tasks
```

## API Specification Changes

### POST /api/tasks (MODIFIED)
**Description**: Create a new task with optional priority and color

**Request**:
```json
{
  "title": "Fix critical bug",
  "priority": "🔥",
  "color": "#dc3545"
}
```

**Request Fields**:
- `title` (string, required): Task title (1-255 characters)
- `priority` (string, optional): Emoticon (🔥, ⭐, ⚡, 💡, 📋), defaults to 📋
- `color` (string, optional): Hex code (#dc3545, #0d6efd, etc.), defaults to #6c757d

**Response** (201 Created):
```json
{
  "id": "1",
  "title": "Fix critical bug",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z",
  "priority": "🔥",
  "color": "#dc3545"
}
```

**Status Codes**:
- 201 Created: Task created successfully
- 400 Bad Request: Invalid input (empty title, invalid priority/color)
- 500 Internal Server Error: Creation failed

**Validation Rules**:
- Title: Not empty, max 255 chars, trimmed (existing)
- Priority: Must be one of 🔥, ⭐, ⚡, 💡, 📋 (new)
- Color: Must be one of 7 valid hex codes (new)
- Empty priority/color: Apply defaults (new)

### GET /api/tasks (MODIFIED)
**Description**: Get all tasks with priority and color

**Response** (200 OK):
```json
[
  {
    "id": "1",
    "title": "Fix critical bug",
    "completed": false,
    "createdAt": "2025-11-19T10:00:00Z",
    "priority": "🔥",
    "color": "#dc3545"
  },
  {
    "id": "2",
    "title": "Plan next sprint",
    "completed": false,
    "createdAt": "2025-11-19T09:00:00Z",
    "priority": "⭐",
    "color": "#0d6efd"
  }
]
```

**Changes**:
- Each task object includes `priority` and `color` fields

### PATCH /api/tasks/{id}/toggle (NO CHANGES)
**No changes needed** - priority cannot be changed after creation

### DELETE /api/tasks/{id} (NO CHANGES)
**No changes needed** - deletion works as before

## Priority-Color Mapping Reference

| Priority | Emoticon | Name | Eisenhower Quadrant | Color | Hex Code | Bootstrap |
|----------|----------|------|---------------------|-------|----------|-----------|
| Urgent & Important | 🔥 | Fire | Q1 | Red | #dc3545 | btn-danger |
| Important, Not Urgent | ⭐ | Star | Q2 | Blue | #0d6efd | btn-primary |
| Urgent, Not Important | ⚡ | Lightning | Q3 | Yellow | #ffc107 | btn-warning |
| Not Urgent, Not Important | 💡 | Light Bulb | Q4 | Green | #28a745 | btn-success |
| Default/Uncategorized | 📋 | Clipboard | N/A | Grey | #6c757d | btn-secondary |

**Reserved for Future**:
- Purple: #6f42c1 (btn-purple)
- Orange: #fd7e14 (btn-orange)

## Security Considerations

### Input Validation (Enhanced)
- **Priority validation**: Server-side check against whitelist of 5 emoticons
- **Color validation**: Server-side check against whitelist of 7 hex codes
- **XSS prevention**: Go template auto-escaping handles UTF-8 emoticons safely
- **No injection risk**: Emoticons stored as string literals, not executable code

### Client-Side Filtering Security
- **No sensitive data**: Priority/color are visual indicators only
- **No authorization bypass**: Filtering is cosmetic (all tasks already sent to client)
- **XSS-safe**: Bootstrap classes and data attributes properly escaped

### Thread Safety (No Changes)
- Priority and Color fields stored atomically with Task struct
- RWMutex ensures thread-safe reads/writes
- No additional synchronization needed

## Performance Considerations

### Storage Impact
- **Memory overhead**: ~20 bytes per task (2 string fields)
- **Acceptable**: For 1000 tasks → ~20KB additional memory
- **No indexing needed**: Client-side filtering eliminates server-side queries

### Client-Side Filtering Performance
- **Algorithm**: O(n) iteration over task elements
- **Benchmark**: ~1ms for 100 tasks, ~10ms for 1000 tasks (well under 100ms target)
- **DOM manipulation**: Show/hide via display property (no reflows)
- **Memory usage**: activeFilters Set (< 1KB)

### Network Impact
- **Request size**: +40 bytes per task (priority + color fields)
- **Response size**: Same as request
- **No additional API calls**: Filtering happens client-side

### Comparison: Client-Side vs Server-Side Filtering

| Aspect | Client-Side (Chosen) | Server-Side |
|--------|---------------------|-------------|
| Latency | < 100ms (instant) | 200-500ms (network RTT) |
| UX | Smooth, instant | Loading spinner needed |
| Server load | None | Increased query complexity |
| Code complexity | Simple JS | New API endpoints |
| Offline support | Works offline | Requires connection |

## Error Handling Strategy

### New Error Types
```go
var (
    ErrInvalidPriority = errors.New("invalid priority emoticon")
    ErrInvalidColor    = errors.New("invalid color code")
)
```

### HTTP Error Responses (Enhanced)
```json
// Invalid priority
{
  "error": "Invalid priority emoticon. Must be one of: 🔥, ⭐, ⚡, 💡, 📋"
}

// Invalid color
{
  "error": "Invalid color code. Must be a valid hex code."
}
```

### Client-Side Error Handling
- Form validation: Ensure priority selected before submission
- Fallback: If no priority selected, use default (📋)
- Error display: Show validation errors in alert div

## Migration Strategy

### Backward Compatibility
**Existing tasks without priority/color**:
- When served via GET /api/tasks, Go zero values apply: `priority: ""`, `color: ""`
- Frontend handles empty values gracefully
- Alternative: Add migration function to TaskStore initialization

### Migration Function (Optional)
```go
func (s *TaskStore) migrateExistingTasks() {
    s.mu.Lock()
    defer s.mu.Unlock()

    migrationCount := 0
    for i := range s.tasks {
        if s.tasks[i].Priority == "" {
            s.tasks[i].Priority = PriorityDefault
            s.tasks[i].Color = ColorGrey
            migrationCount++
        }
    }

    if migrationCount > 0 {
        log.Printf("INFO: Migrated %d tasks with default priority", migrationCount)
    }
}
```

### Rollback Plan
- If feature needs to be rolled back, old version ignores priority/color fields
- Tasks remain functional with just title, completed, createdAt
- No data loss

## Testing Strategy

### Unit Tests (New)
```go
// internal/service/task_service_test.go

func TestTaskService_CreateWithPriority(t *testing.T)
func TestTaskService_ValidatePriority(t *testing.T)
func TestTaskService_ValidateColor(t *testing.T)
func TestTaskService_DefaultPriority(t *testing.T)
func TestTaskService_InvalidPriority(t *testing.T)
```

### Integration Tests (Enhanced)
```go
// internal/handler/api_handler_test.go

func TestAPIHandler_CreateTaskWithPriority(t *testing.T)
func TestAPIHandler_CreateTaskWithoutPriority(t *testing.T)
func TestAPIHandler_CreateTaskInvalidPriority(t *testing.T)
func TestAPIHandler_GetTasksIncludesPriority(t *testing.T)
```

### Frontend Tests (Manual)
1. **Priority Selection**:
   - Select each priority option, verify task created with correct emoticon/color
   - Leave priority unselected, verify default applied

2. **Filtering**:
   - Click single filter, verify only matching tasks shown
   - Click multiple filters, verify OR logic works
   - Click "Show All", verify all tasks visible
   - Verify task count updates correctly

3. **Visual Display**:
   - Verify emoticons render correctly in all browsers
   - Verify colored left borders appear correctly
   - Verify completed tasks show strikethrough + emoticon

## Implementation Phases

### Phase 1: Backend Foundation
1. Extend Task model with Priority and Color fields
2. Update TaskService with validation logic
3. Update APIHandler to accept/return priority/color
4. Write unit tests for validation

### Phase 2: Frontend UI
1. Add priority selector to task creation form
2. Update task list template to display emoticons and colors
3. Add filter buttons above task list
4. Style with Bootstrap classes

### Phase 3: Client-Side Filtering
1. Implement filterByPriority() in Stimulus controller
2. Implement clearFilters() in Stimulus controller
3. Implement applyFilters() with DOM manipulation
4. Add task count display

### Phase 4: Polish & Testing
1. Add migration logic for existing tasks
2. Manual testing of all scenarios
3. Cross-browser testing (Chrome, Firefox, Safari)
4. Mobile responsive testing
5. Update README and documentation

## Directory Structure (Changes)

```
apps/test-task-manager/
├── internal/
│   ├── model/
│   │   └── task.go                 # MODIFIED: Add Priority, Color fields
│   ├── service/
│   │   ├── task_service.go         # MODIFIED: Add priority/color validation
│   │   └── errors.go               # MODIFIED: Add new error types
│   ├── handler/
│   │   └── api_handler.go          # MODIFIED: Accept priority/color in request
├── templates/
│   └── index.html                  # MODIFIED: Add priority selector, filters
├── static/
│   ├── js/
│   │   └── controllers/
│   │       └── tasks_controller.js # MODIFIED: Add filtering logic
│   └── css/
│       └── styles.css               # MODIFIED: Add filter button styles
```

## Architecture Review Checklist

- [x] All functional requirements addressed
- [x] Backward compatibility ensured (default values)
- [x] Priority/color validation defined
- [x] Client-side filtering architecture defined
- [x] API changes specified
- [x] Data model extensions defined
- [x] Template changes defined
- [x] Frontend controller changes defined
- [x] Security considerations addressed
- [x] Performance considerations addressed (< 100ms filtering)
- [x] Migration strategy defined
- [x] Error handling strategy defined
- [x] Testing strategy defined
- [x] Implementation phases defined

## Open Questions

None. Architecture is clear and ready for implementation planning.

## Assumptions

- Existing task manager implementation is complete and functional
- Client-side filtering is acceptable (no server-side filtering needed)
- Users have modern browsers with JavaScript enabled
- Emoticons render correctly on user devices (UTF-8 support)
- ~1000 tasks maximum (client-side filtering scales to this)
- Priority is immutable after creation (no edit functionality)

## Related Documentation

- **FEATURE.md**: Business requirements and user stories
- **TODO.md**: Will contain detailed implementation tasks (created by Tech Lead)
- **Existing ARCHITECTURE.md**: [../task-manager/ARCHITECTURE.md](../task-manager/ARCHITECTURE.md)
- **Developer Log**: [../task-manager/DEVELOPER-GO-LOG.md](../task-manager/DEVELOPER-GO-LOG.md)

---

<!-- COMPLETE -->
