package service

import (
	"errors"
	"testing"

	"gitlab.com/btcdirect-api/test-task-manager/internal/store"
)

func TestTaskService_CreateWithPriority(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	task, err := service.Create("Test task", "🔥", "#dc3545")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Title != "Test task" {
		t.Errorf("expected title 'Test task', got %s", task.Title)
	}
	if task.Priority != "🔥" {
		t.Errorf("expected priority 🔥, got %s", task.Priority)
	}
	if task.Color != "#dc3545" {
		t.Errorf("expected color #dc3545, got %s", task.Color)
	}
}

func TestTaskService_CreateWithDefaults(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	task, err := service.Create("Test task", "", "")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Priority != PriorityDefault {
		t.Errorf("expected default priority %s, got %s", PriorityDefault, task.Priority)
	}
	if task.Color != ColorGrey {
		t.Errorf("expected default color %s, got %s", ColorGrey, task.Color)
	}
}

func TestTaskService_CreateInvalidPriority(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	_, err := service.Create("Test task", "❌", "#dc3545")

	if !errors.Is(err, ErrInvalidPriority) {
		t.Errorf("expected ErrInvalidPriority, got %v", err)
	}
}

func TestTaskService_CreateInvalidColor(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	_, err := service.Create("Test task", "🔥", "#invalid")

	if !errors.Is(err, ErrInvalidColor) {
		t.Errorf("expected ErrInvalidColor, got %v", err)
	}
}

func TestTaskService_CreateEmptyTitle(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	_, err := service.Create("", "🔥", "#dc3545")

	if !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestTaskService_CreateTitleTooLong(t *testing.T) {
	taskStore := store.NewTaskStore()
	service := NewTaskService(taskStore)

	longTitle := make([]byte, 256)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	_, err := service.Create(string(longTitle), "🔥", "#dc3545")

	if !errors.Is(err, ErrTitleTooLong) {
		t.Errorf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		want     bool
	}{
		{"urgent and important", "🔥", true},
		{"important", "⭐", true},
		{"urgent", "⚡", true},
		{"low", "💡", true},
		{"default", "📋", true},
		{"invalid emoticon", "❌", false},
		{"empty string", "", false},
		{"random text", "high", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPriority(tt.priority)
			if got != tt.want {
				t.Errorf("isValidPriority(%q) = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  bool
	}{
		{"red", "#dc3545", true},
		{"blue", "#0d6efd", true},
		{"yellow", "#ffc107", true},
		{"green", "#28a745", true},
		{"purple", "#6f42c1", true},
		{"orange", "#fd7e14", true},
		{"grey", "#6c757d", true},
		{"invalid hex", "#invalid", false},
		{"empty string", "", false},
		{"random text", "red", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidColor(tt.color)
			if got != tt.want {
				t.Errorf("isValidColor(%q) = %v, want %v", tt.color, got, tt.want)
			}
		})
	}
}
