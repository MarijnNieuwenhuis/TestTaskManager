# Simple Task Manager

A lightweight task management web application built with Go and Stimulus.js. This application demonstrates a clean, monolithic web architecture with in-memory storage, perfect for learning and prototyping.

## Features

- **Create Tasks**: Add new tasks with title validation
- **Toggle Completion**: Mark tasks as complete or incomplete
- **Delete Tasks**: Remove tasks with confirmation
- **Real-time Updates**: All interactions via AJAX without page reloads
- **Responsive Design**: Bootstrap 5.3 for mobile and desktop
- **Thread-Safe**: Concurrent access protection with sync.RWMutex
- **Clean Architecture**: Separation of concerns (Model → Store → Service → Handler)

## Technology Stack

**Backend**:
- Go 1.23+ with standard library
- Gorilla Mux for HTTP routing
- html/template for server-side rendering
- In-memory storage (no database required)

**Frontend**:
- Bootstrap 5.3 for styling
- Stimulus.js for progressive enhancement
- Vanilla JavaScript (ES6 modules)

## Getting Started

### Prerequisites

- Go 1.23 or higher
- No database or external dependencies required

### Installation

1. Navigate to the service directory:
```bash
cd apps/test-task-manager
```

2. Configure your `.env` file (optional, defaults provided):
```bash
APP_ENV=dev
HTTP_PORT=8080
LOG_LEVEL=debug
```

### Running Locally

```bash
# Run the application
make run

# Build the binary
make build

# Run the binary
./bin/test-task-manager

# Clean build artifacts
make clean
```

The application will start on `http://localhost:8080`

## Project Structure

```
.
├── cmd/test-task-manager/          # Application entry point
├── internal/
│   ├── app/                        # Application initialization and config
│   ├── model/                      # Data models (Task)
│   ├── store/                      # In-memory storage layer
│   ├── service/                    # Business logic layer
│   ├── handler/                    # HTTP handlers (API + Pages)
│   └── http/
│       ├── handler/                # Legacy health endpoint
│       └── server/                 # Server setup and routing
├── templates/                      # Go HTML templates
│   └── index.html                  # Main task list page
├── static/                         # Static assets
│   ├── css/
│   │   └── styles.css             # Custom styles
│   └── js/
│       ├── app.js                 # Stimulus application bootstrap
│       └── controllers/
│           └── tasks_controller.js # Task interactions controller
├── .env                           # Environment configuration
├── Makefile                       # Build automation
└── README.md                      # This file
```

## Architecture

### Layers

1. **Data Layer**: Task model and in-memory store with thread-safe access
2. **Business Logic Layer**: TaskService with validation (title length, empty check)
3. **HTTP Layer**: Page handler (HTML) and API handler (JSON)
4. **Template Layer**: Go html/template with Bootstrap 5.3
5. **Frontend Layer**: Stimulus.js controllers for dynamic interactions

### API Endpoints

- `GET /` - Main task list page (HTML)
- `GET /health` - Health check endpoint
- `GET /api/tasks` - Get all tasks (JSON)
- `POST /api/tasks` - Create new task (JSON)
- `PATCH /api/tasks/{id}/toggle` - Toggle task completion (JSON)
- `DELETE /api/tasks/{id}` - Delete task (JSON)

### Data Flow

1. User interacts with UI (Stimulus.js)
2. AJAX request to API endpoint
3. API handler validates and calls service
4. Service applies business logic and calls store
5. Store updates in-memory data (thread-safe)
6. Response sent back to frontend
7. UI updates without page reload

## Development

### Task Validation Rules

- Title must not be empty after trimming whitespace
- Title must not exceed 255 characters
- Title is automatically trimmed before saving

### Thread Safety

The TaskStore uses `sync.RWMutex` for concurrent access:
- **Read operations** (GetAll, GetByID): Use RLock for concurrent reads
- **Write operations** (Create, Toggle, Delete): Use Lock for exclusive writes
- All operations return copies to prevent external mutations

### Error Handling

The application uses:
- **Sentinel errors** for expected errors (ErrTaskNotFound, ErrEmptyTitle, ErrTitleTooLong)
- **Error wrapping** with fmt.Errorf and %w for context
- **HTTP status codes**: 200 OK, 201 Created, 400 Bad Request, 404 Not Found, 500 Internal Server Error

## Configuration

Environment variables (configure in `.env`):

- `APP_ENV`: Environment (dev, stage, acc, sandbox, prod) - Default: dev
- `HTTP_PORT`: HTTP server port - Default: 8080
- `LOG_LEVEL`: Logging level (debug, info, warn, error) - Default: info

## Testing

### Manual Testing

1. Start the application: `make run`
2. Open browser: `http://localhost:8080`
3. Test scenarios:
   - Add a task with valid title
   - Try to add empty task (should show error)
   - Try to add very long title (>255 chars, should show error)
   - Toggle task completion (checkbox)
   - Delete a task (with confirmation)
   - Verify all actions work without page reload

### Browser Developer Tools

- Check Network tab for AJAX requests
- Check Console for any JavaScript errors
- Verify response status codes and payloads

## Building

### Local Build
```bash
make build
# Binary created at: bin/test-task-manager
```

### Run Binary
```bash
./bin/test-task-manager -env=dev -port=8080 -loglevel=debug
```

## Features in Detail

### Task Creation
- Form submission via Stimulus.js
- Client-side validation (required field)
- Server-side validation (empty check, length limit)
- Optimistic UI updates
- Error feedback to user

### Task Toggle
- Checkbox interaction via Stimulus.js
- Optimistic UI updates (immediate visual feedback)
- Rollback on server error
- Strikethrough styling for completed tasks

### Task Deletion
- Confirmation dialog before deletion
- Smooth animation on removal
- List automatically updates
- Shows empty state if no tasks remain

### Responsive Design
- Mobile-first Bootstrap 5.3 layout
- Responsive navbar and containers
- Touch-friendly buttons and checkboxes
- Works on desktop, tablet, and mobile

## Troubleshooting

### Port already in use
```bash
# Change port in .env
HTTP_PORT=8081
```

### Templates not found
```bash
# Ensure you run from the correct directory
cd apps/test-task-manager
make run
```

### Static files not loading
- Verify `static/` directory exists
- Check browser console for 404 errors
- Ensure server is running on correct port

## License

Proprietary
