package service

import (
	"rest-todo-api/model/web"

	"golang.org/x/net/context"
)

type TaskService interface {
	Create(ctx context.Context, req web.TaskCreateRequest) (web.TaskResponse, error)
}
