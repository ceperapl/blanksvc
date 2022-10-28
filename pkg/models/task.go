package models

import "time"

type Task struct {
	ID          string     `json:"id" validate:"required,uuid4"`
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description"`
	Deadline    time.Time  `json:"deadline"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
