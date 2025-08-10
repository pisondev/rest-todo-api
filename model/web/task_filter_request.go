package web

import (
	"rest-todo-api/model/domain"
)

type TaskFilterRequest struct {
	Status  *domain.TaskStatus `form:"status"`
	DueDate *string            `form:"due_date"`
}
