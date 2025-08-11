package web

import (
	"rest-todo-api/model/domain"
)

type TaskFilterRequest struct {
	UserID  int                `json:"user_id"`
	Status  *domain.TaskStatus `form:"status"`
	DueDate *string            `form:"due_date"`
}
