package repository

import (
	"rest-todo-api/model/domain"
	"time"
)

type TaskFilter struct {
	Status  *domain.TaskStatus
	DueDate *time.Time
}
