package web

import (
	"rest-todo-api/model/domain"
	"time"
)

type TaskResponse struct {
	ID          int                `json:"id"`
	UserID      int                `json:"user_id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
	DueDate     *time.Time         `json:"due_date"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
