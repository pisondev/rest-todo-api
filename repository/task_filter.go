package repository

import (
	"rest-todo-api/model/domain"
	"time"
)

type TaskFilter struct {
	Status  *domain.TaskStatus
	UserID  int
	DueDate *time.Time
}
