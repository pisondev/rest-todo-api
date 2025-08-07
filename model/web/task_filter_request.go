package web

import (
	"rest-todo-api/model/domain"
	"time"
)

type TaskFilterRequest struct {
	Status  *domain.TaskStatus `form:"status"`
	DueDate *time.Time         `form:"due_date"`
}
