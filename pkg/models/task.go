package models

import "time"

type Task struct {
	ID          string     `json:"id" filtering:"id" validate:"requered"`
	Name        string     `json:"name" filtering:"name" sorting:"name" validate:"requered"`
	Description *string    `json:"description" filtering:"description" sorting:"description"`
	Deadline    string     `json:"deadline" filtering:"deadline" sorting:"deadline" validate:"datetime=2006-01-02"`
	CompletedAt *time.Time `json:"completedAt" filtering:"completed_at" sorting:"completed_at"`
	CreatedAt   time.Time  `json:"createdAt" filtering:"created_at" sorting:"created_at,default:asc"`
	UpdatedAt   *time.Time `json:"updatedAt" filtering:"updated_at" sorting:"updated_at"`
}
