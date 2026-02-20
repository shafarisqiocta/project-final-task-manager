package models

import "time"

type Task struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"project_id"`
	CategoryID  int        `json:"category_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"` // "todo", "in_progress", "done"
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
