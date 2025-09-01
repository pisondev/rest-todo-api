package web

import (
	"rest-todo-api/model/domain"
)

type TaskUpdateRequest struct {
	ID          int                `json:"id"`
	UserID      int                `json:"userId"`
	Title       *string            `validate:"omitempty,min=1" json:"title,omitempty"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
	DueDate     *string            `json:"dueDate"`
}
