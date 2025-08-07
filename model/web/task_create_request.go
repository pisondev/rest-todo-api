package web

import "rest-todo-api/model/domain"

type TaskCreateRequest struct {
	Title       string             `validate:"required,min=1,max=50" json:"title"`
	Description *string            `validate:"max=200" json:"description"`
	Status      *domain.TaskStatus `json:"status"`
}
