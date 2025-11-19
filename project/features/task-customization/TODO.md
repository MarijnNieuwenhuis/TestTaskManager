# Feature: Task Customization with Priority Indicators - Implementation Plan

**Status**: Not Started
**Created**: 2025-11-19
**Go Version**: 1.25
**Architecture**: [ARCHITECTURE.md](./ARCHITECTURE.md)
**Requirements**: [FEATURE.md](./FEATURE.md)

## Overview

This implementation plan adds visual priority indicators (emoticons and colors) to the task manager based on the Eisenhower Matrix methodology, with instant client-side filtering. Users can select from 5 priority levels when creating tasks, and filter the task list by priority using buttons.

## Architecture Summary

**Pattern**: Enhancement to existing Model-Service-Handler architecture
**Components**: Task Model, TaskStore, TaskService, APIHandler, Templates, Stimulus.js Controller
**Complexity**: Medium
**Approach**: Extend existing components with priority/color fields and add client-side filtering logic

### Key Design Decisions (from ARCHITECTURE.md)
- **Backward Compatibility**: Existing tasks automatically receive default priority (📋 Grey)
- **Immutability**: Priority cannot be changed after task creation
- **Client-Side Filtering**: Instant performance (< 100ms target) by filtering in browser
- **Validation**: Server-side validation of priority emoticons and color hex codes
- **5 Priority Levels**: 🔥 Red (Urgent & Important), ⭐ Blue (Important), ⚡ Yellow (Urgent), 💡 Green (Low), 📋 Grey (Default)

## Implementation Phases

### Phase 1: Backend Foundation (4 tasks)
- [ ] Task 1: Extend Task model with Priority and Color fields
- [ ] Task 2: Update TaskStore.Create() to accept priority and color
- [ ] Task 3: Add validation logic and error types to TaskService
- [ ] Task 4: Update TaskService.Create() signature and implementation

### Phase 2: HTTP API Layer (2 tasks)
- [ ] Task 5: Update APIHandler CreateTask request/response structs
- [ ] Task 6: Add error handling for priority/color validation

### Phase 3: Frontend UI (3 tasks)
- [ ] Task 7: Add priority selector to task creation form in templates
- [ ] Task 8: Update task list display with emoticons and colored borders
- [ ] Task 9: Add filter buttons above task list

### Phase 4: Client-Side Filtering (3 tasks)
- [ ] Task 10: Extend Stimulus controller with priority/color handling
- [ ] Task 11: Implement filtering methods (filterByPriority, clearFilters, applyFilters)
- [ ] Task 12: Add task count display logic

### Phase 5: Testing & Validation (3 tasks)
- [ ] Task 13: Write unit tests for validation logic
- [ ] Task 14: Manual testing of all functionality
- [ ] Task 15: Update README documentation

---

## Detailed Task Breakdown

### Task 1: Extend Task Model with Priority and Color Fields

**Phase**: 1 - Backend Foundation
**Dependencies**: None (foundation task)
**Location**: `apps/test-task-manager/internal/model/task.go`
**Estimated Effort**: Small

#### Description
Add two new fields to the Task struct: `Priority` (string containing emoticon) and `Color` (string containing hex code). These fields will store the visual indicators for task prioritization.

#### Acceptance Criteria
- [ ] `Priority` field added to Task struct with `json:"priority"` tag
- [ ] `Color` field added to Task struct with `json:"color"` tag
- [ ] Both fields are string types
- [ ] Fields included in JSON serialization
- [ ] Code compiles without errors

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Data Layer
- Reference: `.claude/refs/go/idiomatic-go.md` - Struct design patterns
- Pattern: Extend existing struct with new fields

#### Files to Create/Modify
- `apps/test-task-manager/internal/model/task.go` - Add Priority and Color fields

#### Go Best Practices
- Use clear field names (Priority, Color)
- Add JSON struct tags for API serialization
- Keep fields exported (capitalized) for JSON marshaling
- Use string type for UTF-8 emoticon support

#### Code Example
```go
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

#### Testing Requirements
- No unit tests needed (data structure only)
- Will be tested via service layer tests

---

### Task 2: Update TaskStore.Create() to Accept Priority and Color

**Phase**: 1 - Backend Foundation
**Dependencies**: Task 1 (requires updated Task struct)
**Location**: `apps/test-task-manager/internal/store/task_store.go`
**Estimated Effort**: Small

#### Description
Modify the `TaskStore.Create()` method signature to accept `title, priority, color` parameters instead of just `title`. Update the method to include priority and color when creating the Task struct.

#### Acceptance Criteria
- [ ] Method signature changed from `Create(title string)` to `Create(title, priority, color string)`
- [ ] Task struct instantiation includes Priority and Color fields
- [ ] All existing functionality preserved (ID, Title, Completed, CreatedAt)
- [ ] Method returns Task with all fields populated
- [ ] Code compiles without errors

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Data Layer → Task Store
- Reference: `.claude/refs/go/best-practices.md` - Method signatures
- Pattern: Extend method parameters

#### Files to Create/Modify
- `apps/test-task-manager/internal/store/task_store.go` - Update Create() method signature and implementation

#### Go Best Practices
- Maintain existing method behavior (thread safety with Lock)
- Keep defer Unlock pattern for safety
- Preserve ID generation logic (strconv.Itoa(s.nextID))
- Return created task with all fields

#### Code Example
```go
// Create adds a new task with priority and color.
func (s *TaskStore) Create(title, priority, color string) model.Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    task := model.Task{
        ID:        strconv.Itoa(s.nextID),
        Title:     title,
        Completed: false,
        CreatedAt: time.Now(),
        Priority:  priority,  // NEW
        Color:     color,     // NEW
    }

    s.tasks = append(s.tasks, task)
    s.nextID++

    return task
}
```

#### Testing Requirements
- No unit tests for this task (will be covered by service tests)
- Store behavior remains the same (just accepts more parameters)

---

### Task 3: Add Validation Logic and Error Types to TaskService

**Phase**: 1 - Backend Foundation
**Dependencies**: Task 1 (requires Task model changes)
**Location**: `apps/test-task-manager/internal/service/errors.go`, `apps/test-task-manager/internal/service/task_service.go`
**Estimated Effort**: Medium

#### Description
Define priority and color constants, create new error types for invalid priority/color, and implement validation functions `isValidPriority()` and `isValidColor()` to ensure only whitelisted emoticons and hex codes are accepted.

#### Acceptance Criteria
- [ ] Priority constants defined (PriorityUrgentImportant, PriorityImportant, PriorityUrgent, PriorityLow, PriorityDefault)
- [ ] Color constants defined (ColorRed, ColorBlue, ColorYellow, ColorGreen, ColorPurple, ColorOrange, ColorGrey)
- [ ] ErrInvalidPriority error created
- [ ] ErrInvalidColor error created
- [ ] isValidPriority() function validates against 5 emoticons
- [ ] isValidColor() function validates against 7 hex codes
- [ ] Functions return true for valid values, false for invalid
- [ ] All error types exported (capitalized)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Business Logic Layer
- Reference: `.claude/refs/go/error-handling.md` - Sentinel errors pattern
- Pattern: Use const for valid values, sentinel errors for validation failures

#### Files to Create/Modify
- `apps/test-task-manager/internal/service/errors.go` - Add ErrInvalidPriority, ErrInvalidColor
- `apps/test-task-manager/internal/service/task_service.go` - Add constants and validation functions

#### Go Best Practices
- Use `const` for immutable values (priority emoticons, color codes)
- Use sentinel errors (errors.New("message")) for validation failures
- Export errors (ErrInvalidPriority) so handlers can check with errors.Is()
- Keep validation functions private (isValidPriority) - internal to service
- Use slice iteration for validation (simple and clear)

#### Code Example
```go
// In errors.go
var (
    ErrEmptyTitle      = errors.New("title cannot be empty")
    ErrTitleTooLong    = errors.New("title cannot exceed 255 characters")
    ErrInvalidPriority = errors.New("invalid priority emoticon")  // NEW
    ErrInvalidColor    = errors.New("invalid color code")         // NEW
)

// In task_service.go
const (
    // Valid priority emoticons
    PriorityUrgentImportant = "🔥"  // Red
    PriorityImportant       = "⭐"  // Blue
    PriorityUrgent          = "⚡"  // Yellow
    PriorityLow             = "💡"  // Green
    PriorityDefault         = "📋"  // Grey

    // Valid color hex codes
    ColorRed    = "#dc3545"
    ColorBlue   = "#0d6efd"
    ColorYellow = "#ffc107"
    ColorGreen  = "#28a745"
    ColorPurple = "#6f42c1"
    ColorOrange = "#fd7e14"
    ColorGrey   = "#6c757d"
)

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

#### Testing Requirements
- Unit tests: `TestIsValidPriority` (test all 5 valid emoticons, test invalid emoticon)
- Unit tests: `TestIsValidColor` (test all 7 valid colors, test invalid color)
- Coverage target: 100% for validation functions

---

### Task 4: Update TaskService.Create() Signature and Implementation

**Phase**: 1 - Backend Foundation
**Dependencies**: Task 2 (requires TaskStore.Create() changes), Task 3 (requires validation functions)
**Location**: `apps/test-task-manager/internal/service/task_service.go`
**Estimated Effort**: Medium

#### Description
Modify `TaskService.Create()` to accept priority and color parameters, apply default values if empty, validate priority and color using the validation functions from Task 3, and call the updated `TaskStore.Create()` method.

#### Acceptance Criteria
- [ ] Method signature changed from `Create(title string)` to `Create(title, priority, color string)`
- [ ] Empty priority defaults to PriorityDefault (📋)
- [ ] Empty color defaults to ColorGrey (#6c757d)
- [ ] Priority validation performed using isValidPriority()
- [ ] Color validation performed using isValidColor()
- [ ] Returns ErrInvalidPriority if priority invalid
- [ ] Returns ErrInvalidColor if color invalid
- [ ] Existing title validation preserved (empty check, length check)
- [ ] Calls TaskStore.Create(title, priority, color)
- [ ] Returns created task and nil error on success

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Business Logic Layer
- Reference: `.claude/refs/go/best-practices.md` - Error handling, validation patterns
- Pattern: Validate all inputs before calling store layer

#### Files to Create/Modify
- `apps/test-task-manager/internal/service/task_service.go` - Update Create() method

#### Go Best Practices
- Check and apply defaults early in function
- Validate in order: title → priority → color (fail fast)
- Use descriptive error returns (ErrInvalidPriority, not generic error)
- Preserve existing behavior (title trimming, validation)
- Return errors directly (no wrapping needed - using sentinel errors)

#### Code Example
```go
// Create creates a new task with validation
func (s *TaskService) Create(title, priority, color string) (model.Task, error) {
    // Validate title (existing logic)
    title = strings.TrimSpace(title)
    if title == "" {
        return model.Task{}, ErrEmptyTitle
    }
    if len(title) > 255 {
        return model.Task{}, ErrTitleTooLong
    }

    // Apply defaults if not provided
    if priority == "" {
        priority = PriorityDefault
    }
    if color == "" {
        color = ColorGrey
    }

    // Validate priority
    if !isValidPriority(priority) {
        return model.Task{}, ErrInvalidPriority
    }

    // Validate color
    if !isValidColor(color) {
        return model.Task{}, ErrInvalidColor
    }

    // Create task with priority and color
    task := s.store.Create(title, priority, color)
    return task, nil
}
```

#### Testing Requirements
- Unit tests: `TestTaskService_CreateWithPriority` (valid priority/color)
- Unit tests: `TestTaskService_CreateWithDefaults` (empty priority/color applies defaults)
- Unit tests: `TestTaskService_CreateInvalidPriority` (returns ErrInvalidPriority)
- Unit tests: `TestTaskService_CreateInvalidColor` (returns ErrInvalidColor)
- Unit tests: `TestTaskService_CreateExistingValidation` (title empty, title too long still work)
- Coverage target: 100% for Create method

---

### Task 5: Update APIHandler CreateTask Request/Response Structs

**Phase**: 2 - HTTP API Layer
**Dependencies**: Task 4 (requires TaskService.Create() changes)
**Location**: `apps/test-task-manager/internal/handler/api_handler.go`
**Estimated Effort**: Small

#### Description
Extend the `CreateTask` request struct to include optional `Priority` and `Color` fields, and update the handler to pass these values to the service layer.

#### Acceptance Criteria
- [ ] Request struct has Priority field with `json:"priority"` tag
- [ ] Request struct has Color field with `json:"color"` tag
- [ ] Both fields are optional (can be empty strings)
- [ ] Handler calls `h.service.Create(req.Title, req.Priority, req.Color)`
- [ ] Response automatically includes priority/color (Task struct has them)
- [ ] Existing behavior preserved (decode JSON, validate, respond)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → HTTP Layer → API Handler
- Reference: `.claude/refs/go/best-practices.md` - HTTP handlers, JSON handling
- Pattern: Extend request struct with optional fields

#### Files to Create/Modify
- `apps/test-task-manager/internal/handler/api_handler.go` - Update CreateTask method

#### Go Best Practices
- Use inline struct definition for request (existing pattern)
- JSON tags match frontend field names (priority, color)
- Let service layer handle defaults (don't set in handler)
- Keep handler thin (just decode, call service, encode response)

#### Code Example
```go
// CreateTask creates a new task from JSON.
func (h *APIHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title    string `json:"title"`
        Priority string `json:"priority"` // NEW: Optional
        Color    string `json:"color"`    // NEW: Optional
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "Invalid request body", "INVALID_INPUT", http.StatusBadRequest)
        return
    }

    task, err := h.service.Create(req.Title, req.Priority, req.Color)  // MODIFIED: Pass priority, color
    if err != nil {
        // Error handling in next task
        ...
    }

    respondJSON(w, task, http.StatusCreated)
}
```

#### Testing Requirements
- Integration tests: `TestAPIHandler_CreateWithPriority` (POST with priority/color)
- Integration tests: `TestAPIHandler_CreateWithoutPriority` (POST without priority/color, defaults applied)
- Coverage target: Handler call path covered

---

### Task 6: Add Error Handling for Priority/Color Validation

**Phase**: 2 - HTTP API Layer
**Dependencies**: Task 5 (requires CreateTask changes)
**Location**: `apps/test-task-manager/internal/handler/api_handler.go`
**Estimated Effort**: Small

#### Description
Add error handling cases in `APIHandler.CreateTask()` to check for `ErrInvalidPriority` and `ErrInvalidColor`, returning appropriate HTTP 400 Bad Request responses with descriptive error messages.

#### Acceptance Criteria
- [ ] Handler checks for ErrInvalidPriority using errors.Is()
- [ ] Handler checks for ErrInvalidColor using errors.Is()
- [ ] Invalid priority returns 400 with "Invalid priority emoticon" message
- [ ] Invalid color returns 400 with "Invalid color code" message
- [ ] Existing error handling preserved (ErrEmptyTitle, ErrTitleTooLong)
- [ ] Error response uses respondError() helper with proper status codes

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → HTTP Layer → API Handler
- Reference: `.claude/refs/go/error-handling.md` - Error checking with errors.Is()
- Pattern: Use switch statement with errors.Is() for error type checking

#### Files to Create/Modify
- `apps/test-task-manager/internal/handler/api_handler.go` - Add error cases to CreateTask

#### Go Best Practices
- Use errors.Is() for sentinel error comparison (not ==)
- Return specific HTTP status codes (400 for validation, 500 for server errors)
- Provide clear error messages for API consumers
- Follow existing error handling pattern in codebase

#### Code Example
```go
task, err := h.service.Create(req.Title, req.Priority, req.Color)
if err != nil {
    if errors.Is(err, service.ErrEmptyTitle) || errors.Is(err, service.ErrTitleTooLong) {
        respondError(w, err.Error(), "INVALID_INPUT", http.StatusBadRequest)
        return
    }
    // NEW: Handle priority validation error
    if errors.Is(err, service.ErrInvalidPriority) {
        respondError(w, "Invalid priority emoticon. Must be one of: 🔥, ⭐, ⚡, 💡, 📋", "INVALID_INPUT", http.StatusBadRequest)
        return
    }
    // NEW: Handle color validation error
    if errors.Is(err, service.ErrInvalidColor) {
        respondError(w, "Invalid color code. Must be a valid hex code.", "INVALID_INPUT", http.StatusBadRequest)
        return
    }
    respondError(w, "Failed to create task", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
    return
}

respondJSON(w, task, http.StatusCreated)
```

#### Testing Requirements
- Integration tests: `TestAPIHandler_CreateInvalidPriority` (returns 400 with error message)
- Integration tests: `TestAPIHandler_CreateInvalidColor` (returns 400 with error message)
- Coverage target: All error paths tested

---

### Task 7: Add Priority Selector to Task Creation Form

**Phase**: 3 - Frontend UI
**Dependencies**: Tasks 1-6 (backend must accept priority/color)
**Location**: `apps/test-task-manager/templates/index.html`
**Estimated Effort**: Medium

#### Description
Add a radio button group below the task input field that allows users to select from 5 priority options. Each option displays an emoticon and description, and has a data-color attribute for the associated hex code.

#### Acceptance Criteria
- [ ] Radio button group added to task creation form
- [ ] 5 options: 🔥 Urgent & Important (Red), ⭐ Important (Blue), ⚡ Urgent (Yellow), 💡 Low (Green), 📋 Default (Grey)
- [ ] Each radio button has value set to emoticon (e.g., value="🔥")
- [ ] Each radio button has data-color attribute (e.g., data-color="#dc3545")
- [ ] Each radio button has data-tasks-target="priorityInput" for Stimulus
- [ ] Default option (📋) is pre-selected with checked attribute
- [ ] Radio buttons styled with Bootstrap btn-check and btn-outline-* classes
- [ ] Button group uses role="group" for accessibility

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Template Layer
- Reference: Bootstrap 5.3 button group documentation
- Pattern: Use Bootstrap button group with radio inputs

#### Files to Create/Modify
- `apps/test-task-manager/templates/index.html` - Add priority selector to form

#### Go Best Practices
- Use Go template syntax correctly: {{range}}, {{if}}, {{.Field}}
- Ensure proper HTML escaping (Go templates auto-escape)

#### Code Example
```html
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
```

#### Testing Requirements
- Manual testing: Verify all 5 radio buttons render correctly
- Manual testing: Verify default option is pre-selected
- Manual testing: Verify only one option can be selected at a time
- Manual testing: Test on Chrome, Firefox, Safari

---

### Task 8: Update Task List Display with Emoticons and Colored Borders

**Phase**: 3 - Frontend UI
**Dependencies**: Tasks 1-6 (backend must return priority/color)
**Location**: `apps/test-task-manager/templates/index.html`
**Estimated Effort**: Small

#### Description
Update the task list item template to display the priority emoticon before the task title, and add a colored left border using inline styles based on the task's color field.

#### Acceptance Criteria
- [ ] Each <li> element has data-priority attribute set to {{.Priority}}
- [ ] Each <li> element has style="border-left: 4px solid {{.Color}}"
- [ ] Task label displays emoticon before title: <span>{{.Priority}}</span>{{.Title}}
- [ ] Emoticon wrapped in <span> with spacing class (me-2)
- [ ] Existing checkbox and delete button preserved
- [ ] Completed tasks still show strikethrough with emoticon visible

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Template Layer
- Reference: Go template documentation for field access ({{.Field}})
- Pattern: Use inline styles for dynamic colors, data attributes for filtering

#### Files to Create/Modify
- `apps/test-task-manager/templates/index.html` - Update task list item template

#### Go Best Practices
- Use Go template syntax: {{.Priority}}, {{.Color}}
- Go templates auto-escape HTML (emoticons are safe UTF-8)
- Use data-* attributes for JavaScript access

#### Code Example
```html
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
```

#### Testing Requirements
- Manual testing: Verify emoticons render correctly in all browsers
- Manual testing: Verify colored left borders appear with correct colors
- Manual testing: Verify completed tasks show strikethrough + emoticon
- Manual testing: Verify empty list shows "No tasks" message

---

### Task 9: Add Filter Buttons Above Task List

**Phase**: 3 - Frontend UI
**Dependencies**: Task 8 (requires data-priority attributes on list items)
**Location**: `apps/test-task-manager/templates/index.html`
**Estimated Effort**: Small

#### Description
Add a button group above the task list containing filter buttons for each priority level plus a "Show All" button. Add a task count display that shows the current number of visible tasks.

#### Acceptance Criteria
- [ ] Button group added above task list with 6 buttons
- [ ] "Show All" button with data-action="click->tasks#clearFilters"
- [ ] 5 filter buttons (🔥, ⭐, ⚡, 💡, 📋) each with data-action="click->tasks#filterByPriority"
- [ ] Each filter button has data-priority attribute matching emoticon
- [ ] Filter buttons styled with btn-sm and btn-outline-* matching priority colors
- [ ] Task count span with data-tasks-target="taskCount"
- [ ] Initial count displays {{len .Tasks}} tasks
- [ ] Buttons use appropriate Bootstrap color classes

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Template Layer
- Reference: Bootstrap 5.3 button group, sizing utilities
- Pattern: Data attributes for Stimulus controller actions

#### Files to Create/Modify
- `apps/test-task-manager/templates/index.html` - Add filter buttons above task list

#### Go Best Practices
- Use Go template function: {{len .Tasks}} for initial count
- Use consistent naming for data attributes

#### Code Example
```html
<!-- NEW: Filter Buttons -->
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

<!-- Task List -->
<ul class="list-group" data-tasks-target="list">
    ...
</ul>
```

#### Testing Requirements
- Manual testing: Verify all 6 buttons render correctly
- Manual testing: Verify emoticons display in filter buttons
- Manual testing: Verify task count shows correct initial value
- Manual testing: Verify button colors match priority colors

---

### Task 10: Extend Stimulus Controller with Priority/Color Handling

**Phase**: 4 - Client-Side Filtering
**Dependencies**: Task 7 (requires priority selector in form), Task 9 (requires filter buttons)
**Location**: `apps/test-task-manager/static/js/controllers/tasks_controller.js`
**Estimated Effort**: Medium

#### Description
Update the Stimulus controller's `create()` method to read the selected priority and color from the radio button group, and include them in the POST request to the API. Add new targets for priorityInput and taskCount.

#### Acceptance Criteria
- [ ] Controller targets array includes "priorityInput" and "taskCount"
- [ ] create() method finds checked radio button from priorityInputTargets
- [ ] create() method extracts priority from input.value
- [ ] create() method extracts color from input.dataset.color
- [ ] create() method defaults to "📋" and "#6c757d" if no selection
- [ ] Priority and color included in JSON.stringify({ title, priority, color })
- [ ] Existing create() behavior preserved (validation, error handling, reload)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Frontend Layer
- Reference: Stimulus.js 3.2+ documentation - targets, data attributes
- Pattern: Read form state, include in API request

#### Files to Create/Modify
- `apps/test-task-manager/static/js/controllers/tasks_controller.js` - Update create() method and targets

#### Go Best Practices
- N/A (JavaScript file)

#### Code Example
```javascript
export default class extends Controller {
    static targets = ["input", "error", "list", "label", "priorityInput", "taskCount"]  // MODIFIED: Added priorityInput, taskCount

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
                body: JSON.stringify({ title, priority, color }),  // MODIFIED: Include priority, color
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

    // Existing methods: toggle, delete, showError, hideError...
}
```

#### Testing Requirements
- Manual testing: Create task with each priority, verify API receives priority/color
- Manual testing: Create task without selecting priority, verify defaults applied
- Manual testing: Check browser console for errors
- Manual testing: Verify task appears with correct emoticon and color

---

### Task 11: Implement Filtering Methods

**Phase**: 4 - Client-Side Filtering
**Dependencies**: Task 10 (requires controller structure)
**Location**: `apps/test-task-manager/static/js/controllers/tasks_controller.js`
**Estimated Effort**: Medium

#### Description
Implement three new methods in the Stimulus controller: `filterByPriority()` to toggle filter on/off, `clearFilters()` to reset all filters, and `applyFilters()` to show/hide tasks based on active filters. Use a Set to track active filters.

#### Acceptance Criteria
- [ ] activeFilters Set initialized as class property
- [ ] filterByPriority() reads data-priority from clicked button
- [ ] filterByPriority() toggles priority in activeFilters Set
- [ ] filterByPriority() adds/removes "active" class on button
- [ ] filterByPriority() calls applyFilters()
- [ ] clearFilters() clears activeFilters Set
- [ ] clearFilters() removes "active" class from all filter buttons
- [ ] clearFilters() calls applyFilters()
- [ ] applyFilters() iterates all <li> elements with data-task-id
- [ ] applyFilters() shows task if no filters active OR priority in activeFilters
- [ ] applyFilters() hides task if filters active AND priority not in activeFilters
- [ ] applyFilters() uses display style property (show: "", hide: "none")

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Frontend Layer
- Reference: JavaScript Set documentation, DOM manipulation
- Pattern: Toggle pattern for filters, OR logic for multiple filters

#### Files to Create/Modify
- `apps/test-task-manager/static/js/controllers/tasks_controller.js` - Add filtering methods

#### Go Best Practices
- N/A (JavaScript file)

#### Code Example
```javascript
export default class extends Controller {
    static targets = ["input", "error", "list", "label", "priorityInput", "taskCount"]

    // NEW: Track active filters
    activeFilters = new Set()

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

        tasks.forEach(task => {
            const taskPriority = task.dataset.priority

            // Show if no filters active OR priority matches any active filter
            if (this.activeFilters.size === 0 || this.activeFilters.has(taskPriority)) {
                task.style.display = ""
            } else {
                task.style.display = "none"
            }
        })

        // Task count update in next task
    }

    // Existing methods...
}
```

#### Testing Requirements
- Manual testing: Click 🔥 filter, verify only 🔥 tasks shown
- Manual testing: Click 🔥 again, verify all tasks shown (toggle off)
- Manual testing: Click 🔥 and ⭐, verify both shown (OR logic)
- Manual testing: Click "Show All", verify all filters cleared
- Manual testing: Verify filter buttons show active state when clicked
- Manual testing: Test with empty task list (no errors)

---

### Task 12: Add Task Count Display Logic

**Phase**: 4 - Client-Side Filtering
**Dependencies**: Task 11 (requires applyFilters() method)
**Location**: `apps/test-task-manager/static/js/controllers/tasks_controller.js`
**Estimated Effort**: Small

#### Description
Update the `applyFilters()` method to count visible tasks and update the task count display element with text like "Showing 3 of 10 tasks" when filters are active, or "Showing 10 tasks" when no filters are active.

#### Acceptance Criteria
- [ ] applyFilters() counts visible tasks (visibleCount variable)
- [ ] applyFilters() gets total count (tasks.length)
- [ ] When filters active: displays "Showing X of Y tasks"
- [ ] When no filters: displays "Showing X tasks"
- [ ] Task count updates immediately when filters change
- [ ] Checks hasTaskCountTarget before updating (safe if target missing)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Component Design → Frontend Layer
- Reference: Stimulus.js targets documentation
- Pattern: Update DOM text content based on state

#### Files to Create/Modify
- `apps/test-task-manager/static/js/controllers/tasks_controller.js` - Update applyFilters() method

#### Go Best Practices
- N/A (JavaScript file)

#### Code Example
```javascript
// NEW: Apply active filters to task list
applyFilters() {
    const tasks = this.listTarget.querySelectorAll('[data-task-id]')
    let visibleCount = 0  // NEW: Track visible count

    tasks.forEach(task => {
        const taskPriority = task.dataset.priority

        // Show if no filters active OR priority matches any active filter
        if (this.activeFilters.size === 0 || this.activeFilters.has(taskPriority)) {
            task.style.display = ""
            visibleCount++  // NEW: Increment visible count
        } else {
            task.style.display = "none"
        }
    })

    // NEW: Update task count
    if (this.hasTaskCountTarget) {
        const totalCount = tasks.length
        const countText = this.activeFilters.size > 0
            ? `Showing ${visibleCount} of ${totalCount} tasks`
            : `Showing ${totalCount} tasks`
        this.taskCountTarget.textContent = countText
    }
}
```

#### Testing Requirements
- Manual testing: Verify task count shows "Showing N tasks" initially
- Manual testing: Click filter, verify count shows "Showing X of Y tasks"
- Manual testing: Clear filters, verify count shows "Showing Y tasks"
- Manual testing: Test with 0 tasks (shows "Showing 0 tasks")
- Manual testing: Test with all tasks filtered out (shows "Showing 0 of N tasks")

---

### Task 13: Write Unit Tests for Validation Logic

**Phase**: 5 - Testing & Validation
**Dependencies**: Task 3 (requires validation functions), Task 4 (requires TaskService.Create())
**Location**: `apps/test-task-manager/internal/service/task_service_test.go`
**Estimated Effort**: Medium

#### Description
Write comprehensive unit tests for the priority and color validation logic in TaskService. Test all valid values, invalid values, default application, and error returns.

#### Acceptance Criteria
- [ ] TestIsValidPriority_ValidEmoticons tests all 5 valid emoticons return true
- [ ] TestIsValidPriority_InvalidEmoticon tests invalid emoticon returns false
- [ ] TestIsValidColor_ValidColors tests all 7 valid colors return true
- [ ] TestIsValidColor_InvalidColor tests invalid color returns false
- [ ] TestTaskService_CreateWithPriority tests create with valid priority/color
- [ ] TestTaskService_CreateWithDefaults tests empty priority/color applies defaults
- [ ] TestTaskService_CreateInvalidPriority tests returns ErrInvalidPriority
- [ ] TestTaskService_CreateInvalidColor tests returns ErrInvalidColor
- [ ] TestTaskService_CreateExistingValidation tests title validation still works
- [ ] All tests pass with `go test ./internal/service/...`
- [ ] Coverage > 90% for validation code

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Testing Strategy
- Reference: `.claude/refs/go/testing-practices.md` - Table-driven tests, error checking
- Pattern: Use table-driven tests for multiple valid/invalid values

#### Files to Create/Modify
- `apps/test-task-manager/internal/service/task_service_test.go` - Add test functions

#### Go Best Practices
- Use table-driven tests for testing multiple inputs
- Use t.Run() for subtests with descriptive names
- Use errors.Is() to check for sentinel errors
- Test both success and failure cases
- Mock dependencies (use real TaskStore - it's simple in-memory)

#### Code Example
```go
func TestTaskService_CreateWithPriority(t *testing.T) {
    store := store.NewTaskStore()
    service := NewTaskService(store)

    task, err := service.Create("Test task", "🔥", "#dc3545")

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if task.Priority != "🔥" {
        t.Errorf("expected priority 🔥, got %s", task.Priority)
    }
    if task.Color != "#dc3545" {
        t.Errorf("expected color #dc3545, got %s", task.Color)
    }
}

func TestTaskService_CreateWithDefaults(t *testing.T) {
    store := store.NewTaskStore()
    service := NewTaskService(store)

    task, err := service.Create("Test task", "", "")

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if task.Priority != PriorityDefault {
        t.Errorf("expected default priority, got %s", task.Priority)
    }
    if task.Color != ColorGrey {
        t.Errorf("expected default color, got %s", task.Color)
    }
}

func TestTaskService_CreateInvalidPriority(t *testing.T) {
    store := store.NewTaskStore()
    service := NewTaskService(store)

    _, err := service.Create("Test task", "❌", "#dc3545")

    if !errors.Is(err, ErrInvalidPriority) {
        t.Errorf("expected ErrInvalidPriority, got %v", err)
    }
}

// Additional tests: TestIsValidPriority, TestIsValidColor, TestTaskService_CreateInvalidColor, etc.
```

#### Testing Requirements
- Unit tests: All test functions pass
- Coverage: Run `go test -cover ./internal/service/...` - should be > 90%
- Edge cases: Empty strings, invalid values, all valid values

---

### Task 14: Manual Testing of All Functionality

**Phase**: 5 - Testing & Validation
**Dependencies**: Tasks 1-12 (requires full implementation)
**Location**: N/A (manual testing)
**Estimated Effort**: Medium

#### Description
Perform comprehensive manual testing of the entire feature in a web browser. Test all user stories, acceptance criteria, edge cases, and cross-browser compatibility.

#### Acceptance Criteria
- [ ] All user stories from FEATURE.md tested and passing
- [ ] Priority selection works for all 5 options
- [ ] Tasks display with correct emoticons and colors
- [ ] Filtering works for each priority individually
- [ ] Multiple filters can be combined (OR logic)
- [ ] "Show All" clears all filters
- [ ] Task count updates correctly
- [ ] Filter buttons show active state
- [ ] Create task without priority defaults to 📋 Grey
- [ ] Completed tasks show strikethrough with visible emoticon
- [ ] All functionality tested on Chrome, Firefox, Safari
- [ ] Mobile responsive design works (test on phone/tablet or browser dev tools)
- [ ] No JavaScript errors in browser console
- [ ] API returns correct error messages for invalid priority/color

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Success Criteria
- Reference: FEATURE.md - User Stories, Acceptance Criteria
- Pattern: Follow acceptance criteria as test checklist

#### Files to Create/Modify
- N/A (manual testing only)

#### Go Best Practices
- N/A (manual testing)

#### Testing Checklist

**Priority Selection**:
- [ ] Select 🔥, create task, verify task shows 🔥 with red border
- [ ] Select ⭐, create task, verify task shows ⭐ with blue border
- [ ] Select ⚡, create task, verify task shows ⚡ with yellow border
- [ ] Select 💡, create task, verify task shows 💡 with green border
- [ ] Select 📋, create task, verify task shows 📋 with grey border
- [ ] Don't select anything (default), create task, verify shows 📋 grey

**Filtering**:
- [ ] Click 🔥 filter, verify only 🔥 tasks shown
- [ ] Click 🔥 again, verify all tasks shown (toggle off)
- [ ] Click 🔥 and ⭐, verify both shown, others hidden
- [ ] Click "Show All", verify all filters cleared and all tasks shown
- [ ] Filter with no matching tasks, verify "Showing 0 of N tasks"
- [ ] Filter buttons show active state (highlighted when clicked)

**Task Count**:
- [ ] Initial load shows "Showing N tasks"
- [ ] Apply filter shows "Showing X of N tasks"
- [ ] Clear filter shows "Showing N tasks"

**Edge Cases**:
- [ ] Create task with empty title (should show error)
- [ ] Create task with very long title > 255 chars (should show error)
- [ ] Toggle completed status preserves emoticon and color
- [ ] Delete task removes it from list
- [ ] Page reload preserves tasks (in-memory store)

**Cross-Browser**:
- [ ] Test all functionality in Chrome
- [ ] Test all functionality in Firefox
- [ ] Test all functionality in Safari
- [ ] Verify emoticons render correctly in all browsers
- [ ] Check for any JavaScript errors in console

**Mobile/Responsive**:
- [ ] Priority selector works on touch screens
- [ ] Filter buttons wrap on small screens
- [ ] Emoticons display correctly on mobile
- [ ] All interactions work with touch

**API Testing** (use browser dev tools Network tab):
- [ ] POST /api/tasks includes priority and color in request
- [ ] POST /api/tasks returns task with priority and color
- [ ] Invalid priority returns 400 error with message
- [ ] Invalid color returns 400 error with message

#### Testing Requirements
- All items in checklist completed and passing
- Document any bugs found and fixed
- Verify performance: filtering < 100ms (feels instant)

---

### Task 15: Update README Documentation

**Phase**: 5 - Testing & Validation
**Dependencies**: Tasks 1-14 (requires full implementation and testing)
**Location**: `apps/test-task-manager/README.md`
**Estimated Effort**: Small

#### Description
Update the README.md to document the new priority indicator feature, including how to use priority selection, how filtering works, and the Eisenhower Matrix mapping.

#### Acceptance Criteria
- [ ] Section added: "Priority Indicators" or "Task Prioritization"
- [ ] Documents 5 priority levels with emoticon/color/meaning
- [ ] Explains Eisenhower Matrix methodology
- [ ] Describes how to select priority when creating tasks
- [ ] Describes how to filter tasks by priority
- [ ] Notes that priority is immutable after creation
- [ ] API documentation updated with priority/color fields
- [ ] Architecture section updated (if exists)

#### Implementation Notes
- Architecture: See ARCHITECTURE.md - Priority-Color Mapping Reference
- Reference: FEATURE.md - Priority-Color Mapping table
- Pattern: Add new section after existing features

#### Files to Create/Modify
- `apps/test-task-manager/README.md` - Add priority feature documentation

#### Go Best Practices
- N/A (documentation only)

#### Code Example
```markdown
## Features

- Create, view, toggle, and delete tasks
- **Priority Indicators**: Categorize tasks using the Eisenhower Matrix
- Client-side filtering by priority
- Responsive Bootstrap UI
- Thread-safe in-memory storage

## Priority Indicators

Tasks can be assigned priority levels based on the [Eisenhower Matrix](https://en.wikipedia.org/wiki/Time_management#The_Eisenhower_Method), which categorizes tasks by urgency and importance:

| Priority | Emoticon | Color | Meaning |
|----------|----------|-------|---------|
| Urgent & Important | 🔥 | Red | Do first - critical tasks |
| Important, Not Urgent | ⭐ | Blue | Schedule - important but not immediate |
| Urgent, Not Important | ⚡ | Yellow | Delegate - quick tasks |
| Not Urgent, Not Important | 💡 | Green | Low priority - do later |
| Default/Uncategorized | 📋 | Grey | Standard task |

### Using Priorities

**Selecting Priority**:
1. Type your task title
2. Select a priority option (defaults to 📋 Grey if not selected)
3. Click "Add Task"

**Filtering Tasks**:
- Click filter buttons above the task list to show only tasks with specific priorities
- Click multiple filters to combine them (shows tasks matching ANY selected priority)
- Click "Show All" to clear all filters

**Note**: Priority cannot be changed after task creation. To change priority, delete and recreate the task.

## API Documentation

### POST /api/tasks

**Request**:
```json
{
  "title": "Task title",
  "priority": "🔥",       // Optional: defaults to 📋
  "color": "#dc3545"      // Optional: defaults to #6c757d
}
```

**Response**:
```json
{
  "id": "1",
  "title": "Task title",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z",
  "priority": "🔥",
  "color": "#dc3545"
}
```

### GET /api/tasks

**Response**: Array of tasks (each includes priority and color fields)

*(Update existing API documentation sections similarly)*
```

#### Testing Requirements
- Review: Read updated README for clarity and accuracy
- Verify: All links work (if any added)
- Verify: Code examples are correct
- Verify: Markdown renders correctly (check in GitHub or markdown preview)

---

## Go Best Practices to Follow

### Code Organization (from `.claude/refs/go/best-practices.md`)
- Keep packages focused (model, store, service, handler separation)
- Use clear package names (no "utils" or "helpers")
- Export only what's necessary (validation functions stay private)

### Idiomatic Go (from `.claude/refs/go/idiomatic-go.md`)
- Use short variable names in small scopes (p for priority in validation)
- Accept interfaces, return structs (TaskService returns model.Task)
- Use constants for fixed values (priority emoticons, color codes)
- Keep error messages lowercase, no punctuation (idiomatic error format)

### Error Handling (from `.claude/refs/go/error-handling.md`)
- Use sentinel errors for validation failures (errors.New())
- Export error variables for use with errors.Is()
- Check errors with errors.Is(), not == comparison
- Wrap errors only when adding context (validation errors don't need wrapping)

### Testing (from `.claude/refs/go/testing-practices.md`)
- Use table-driven tests for multiple similar cases
- Use t.Run() for subtests with descriptive names
- Test both success and error paths
- Use errors.Is() in tests for error checking
- Aim for >80% coverage on business logic

### Concurrency (from `.claude/refs/go/concurrency-patterns.md`)
- TaskStore already uses sync.RWMutex correctly
- No new goroutines needed (synchronous HTTP handlers)
- Maintain existing lock/unlock patterns
- Use defer for unlock (already in place)

### Design Patterns (from `.claude/refs/go/design-patterns.md`)
- Repository pattern: TaskStore abstracts storage
- Service pattern: TaskService handles business logic
- Dependency injection: Handlers receive service via constructor
- Validation pattern: Fail fast with clear errors

## Validation Checklist

### Requirements Coverage
- [x] All functional requirements from FEATURE.md addressed
- [x] All non-functional requirements addressed (performance, usability, maintainability)
- [x] All user stories have corresponding implementation tasks
- [x] All acceptance criteria covered in tasks

### Architecture Adherence
- [x] All components from ARCHITECTURE.md have implementation tasks
- [x] Task model extended per architecture (Priority, Color fields)
- [x] Service layer validation implemented per architecture
- [x] API handler changes per architecture
- [x] Template changes per architecture
- [x] Stimulus controller changes per architecture

### Go Best Practices
- [x] References `.claude/refs/go/` resources in each task
- [x] Follows Go idioms (constants, sentinel errors, validation)
- [x] Proper error handling (errors.Is(), sentinel errors)
- [x] Testing strategy defined
- [x] Concurrency handled (existing RWMutex preserved)

### Task Quality
- [x] All tasks are specific and actionable
- [x] Acceptance criteria defined for each task
- [x] Dependencies clearly stated
- [x] Files to create/modify listed
- [x] Effort estimated (Small/Medium/Large)
- [x] Code examples provided where helpful

### Technical Soundness
- [x] Architecture translated correctly to tasks
- [x] No design flaws introduced
- [x] Performance considered (client-side filtering)
- [x] Security addressed (validation, XSS prevention)
- [x] Testing comprehensive (unit + manual)

### Completeness
- [x] Backend tasks complete (model, store, service, handler)
- [x] Frontend tasks complete (templates, Stimulus controller)
- [x] Testing tasks defined (unit tests, manual testing)
- [x] Documentation task included (README update)
- [x] All phases have clear boundaries

## Success Criteria (from FEATURE.md)

### Definition of Done
- [ ] All functional requirements implemented (Tasks 1-12)
- [ ] Task model extended with Priority and Color fields (Task 1)
- [ ] API accepts and returns priority/color in JSON (Tasks 5-6)
- [ ] Service layer validates priority/color values (Tasks 3-4)
- [ ] Templates render emoticons and color styling (Tasks 7-8)
- [ ] Stimulus controller handles priority selection (Task 10)
- [ ] Client-side filtering implemented (Task 11)
- [ ] Filter buttons rendered and functional (Task 9)
- [ ] Default priority applied to new tasks (Task 4)
- [ ] All acceptance criteria met (verified in Task 14)
- [ ] Code reviewed and approved (Team review after implementation)
- [ ] Manual testing completed (Task 14)
- [ ] Documentation updated (Task 15)
- [ ] Deployed and verified working (After implementation)

### Performance Targets
- Filtering response time < 100ms ✓ (client-side filtering in Task 11)
- Visual distinction immediately clear ✓ (emoticons + colors in Tasks 7-8)
- All existing tasks receive default categorization ✓ (defaults in Task 4)

### Quality Targets
- Unit test coverage > 90% for validation logic (Task 13)
- Cross-browser compatibility (Chrome, Firefox, Safari) (Task 14)
- Mobile responsive design works (Task 14)
- No JavaScript errors in console (Task 14)

## Implementation Timeline

**Estimated Total Effort**:
- Phase 1 (Backend Foundation): 4 tasks × Medium = ~4-6 hours
- Phase 2 (HTTP API Layer): 2 tasks × Small = ~1-2 hours
- Phase 3 (Frontend UI): 3 tasks × Medium = ~3-4 hours
- Phase 4 (Client-Side Filtering): 3 tasks × Medium = ~3-4 hours
- Phase 5 (Testing & Validation): 3 tasks × Medium = ~4-6 hours

**Total**: ~15-22 hours for complete implementation

**Suggested Approach**:
1. Implement Phases 1-2 together (backend complete)
2. Implement Phases 3-4 together (frontend complete)
3. Phase 5 (testing and documentation)

## Notes for Developer

### Working Directory
All file paths are relative to: `apps/test-task-manager/`

Example:
- Model: `internal/model/task.go`
- Service: `internal/service/task_service.go`
- Templates: `templates/index.html`
- Frontend: `static/js/controllers/tasks_controller.js`

### Development Workflow
1. Create feature branch: `git checkout -b feature/task-customization`
2. Implement tasks sequentially (1 → 15)
3. Commit after each phase completes
4. Run tests frequently: `go test ./internal/...`
5. Test in browser after frontend tasks
6. Create PR when all tasks complete

### Testing Commands
```bash
# Run all tests
go test ./internal/...

# Run tests with coverage
go test -cover ./internal/...

# Run specific test
go test ./internal/service -run TestTaskService_CreateWithPriority

# Build project
go build ./cmd/test-task-manager

# Run project
go run ./cmd/test-task-manager/main.go
```

### Browser Testing
1. Run application: `make run` or `go run ./cmd/test-task-manager/main.go`
2. Open: `http://localhost:8080`
3. Open browser dev tools (F12)
4. Check Network tab for API calls
5. Check Console tab for JavaScript errors

### Debugging Tips
- Use `fmt.Printf()` for debugging Go code
- Use `console.log()` for debugging JavaScript
- Check browser Network tab for API request/response
- Verify JSON structure matches expected format
- Test emoticons render correctly (UTF-8 support)

### Common Pitfalls to Avoid
- Don't forget to add `json` struct tags to new fields
- Don't forget to use `errors.Is()` instead of `==` for error checking
- Don't forget to add Stimulus targets to static targets array
- Don't forget to handle empty priority/color (apply defaults)
- Don't forget checked attribute on default priority radio button
- Don't forget data-priority attribute on task list items (needed for filtering)

### When You're Stuck
- Re-read ARCHITECTURE.md for design decisions
- Re-read FEATURE.md for requirements
- Check existing code patterns (how task creation currently works)
- Reference `.claude/refs/go/` documentation for Go best practices
- Test incrementally (don't implement everything before testing)

---

## Related Documentation

- **FEATURE.md**: [Business requirements and user stories](./FEATURE.md)
- **ARCHITECTURE.md**: [Technical architecture and design](./ARCHITECTURE.md)
- **Go Best Practices**: `../.claude/refs/go/best-practices.md`
- **Go Error Handling**: `../.claude/refs/go/error-handling.md`
- **Go Testing**: `../.claude/refs/go/testing-practices.md`
- **Idiomatic Go**: `../.claude/refs/go/idiomatic-go.md`
- **Task Manager Feature**: `../task-manager/` (reference existing implementation)

---

<!-- COMPLETE -->
