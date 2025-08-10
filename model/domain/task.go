package domain

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int
	UserID      int
	Title       string
	Description *string
	Status      *TaskStatus
	DueDate     sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusDone    TaskStatus = "done"
)
