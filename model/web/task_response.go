package web

import (
	"rest-todo-api/model/domain"
	"time"
)

type TaskResponse struct {
	ID          int                `json:"id"`
	UserID      int                `json:"userId"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
	DueDate     *time.Time         `json:"dueDate"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}
