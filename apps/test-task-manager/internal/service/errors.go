package service

import "errors"

var (
	// ErrEmptyTitle is returned when a task title is empty.
	ErrEmptyTitle = errors.New("task title cannot be empty")
	// ErrTitleTooLong is returned when a task title exceeds 255 characters.
	ErrTitleTooLong = errors.New("task title cannot exceed 255 characters")
)
