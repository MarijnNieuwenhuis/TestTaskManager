// Package model defines the data models for the task manager.
package model

import "time"

// Task represents a single task item in the task manager.
type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}
