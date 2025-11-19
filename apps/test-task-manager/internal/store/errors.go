package store

import "errors"

// ErrTaskNotFound is returned when a task with the given ID doesn't exist.
var ErrTaskNotFound = errors.New("task not found")
