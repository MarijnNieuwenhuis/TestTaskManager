// Package service implements business logic for the task manager.
package service

import (
	"fmt"
	"strings"

	"gitlab.com/btcdirect-api/test-task-manager/internal/model"
	"gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

// TaskService handles business logic for tasks.
type TaskService struct {
	store *store.TaskStore
}

// NewTaskService creates a new TaskService.
func NewTaskService(store *store.TaskStore) *TaskService {
	return &TaskService{store: store}
}

// GetAll retrieves all tasks.
func (s *TaskService) GetAll() []model.Task {
	return s.store.GetAll()
}

// Create creates a new task with validation.
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

// Toggle toggles task completion status.
func (s *TaskService) Toggle(id string) (model.Task, error) {
	task, err := s.store.Toggle(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to toggle task: %w", err)
	}
	return task, nil
}

// Delete removes a task.
func (s *TaskService) Delete(id string) error {
	if err := s.store.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
