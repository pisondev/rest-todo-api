package web

import (
	"rest-todo-api/model/domain"
)

type TaskCreateRequest struct {
	UserID      int                `json:"user_id"`
	Title       string             `validate:"required,min=1,max=50" json:"title"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
	DueDate     *string            `json:"due_date"`
}
