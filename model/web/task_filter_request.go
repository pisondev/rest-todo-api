package web

import (
	"rest-todo-api/model/domain"
)

type TaskFilterRequest struct {
	UserID  int                `json:"userId"`
	Status  *domain.TaskStatus `form:"status"`
	DueDate *string            `form:"dueDate"`
}
