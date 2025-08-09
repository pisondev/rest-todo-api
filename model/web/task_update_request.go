package web

import "rest-todo-api/model/domain"

type TaskUpdateRequest struct {
	ID          *int               `json:"id"`
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
}
