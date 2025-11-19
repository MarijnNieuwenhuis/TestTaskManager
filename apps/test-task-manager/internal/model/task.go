// Package model defines the data models for the task manager.
package model

import "time"

// Task represents a single task item in the task manager with priority indicators.
type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	Priority  string    `json:"priority"` // Emoticon representing priority (🔥, ⭐, ⚡, 💡, 📋)
	Color     string    `json:"color"`    // Hex color code for visual display
}
