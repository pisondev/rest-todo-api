package domain

import (
	"time"
)

type Task struct {
	ID          int
	UserID      int
	Title       string
	Description *string
	Status      *TaskStatus
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusDone    TaskStatus = "done"
)
