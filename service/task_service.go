package service

import (
	"rest-todo-api/model/web"

	"golang.org/x/net/context"
)

type TaskService interface {
	Create(ctx context.Context, req web.TaskCreateRequest) (web.TaskResponse, error)
	FindTasks(ctx context.Context, req web.TaskFilterRequest) ([]web.TaskResponse, error)
	FindByID(ctx context.Context, taskID int) (web.TaskResponse, error)
	Update(ctx context.Context, req web.TaskUpdateRequest) (web.TaskResponse, error)
}
