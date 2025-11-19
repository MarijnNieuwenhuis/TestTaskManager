// Package store provides thread-safe in-memory task storage.
package store

import (
	"strconv"
	"sync"
	"time"

	"gitlab.com/btcdirect-api/test-task-manager/internal/model"
)

// TaskStore provides thread-safe in-memory task storage.
type TaskStore struct {
	tasks  []model.Task
	nextID int
	mu     sync.RWMutex
}

// NewTaskStore creates a new TaskStore.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make([]model.Task, 0),
		nextID: 1,
	}
}

// GetAll returns all tasks.
func (s *TaskStore) GetAll() []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	tasksCopy := make([]model.Task, len(s.tasks))
	copy(tasksCopy, s.tasks)
	return tasksCopy
}

// GetByID returns a task by ID.
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

// Create adds a new task.
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

// Toggle changes completion status.
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

// Delete removes a task.
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
