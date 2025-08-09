package service

import (
	"context"
	"database/sql"
	"rest-todo-api/helper"
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
	"rest-todo-api/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type TaskServiceImpl struct {
	TaskRepository repository.TaskRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewTaskService(taskRepository repository.TaskRepository, DB *sql.DB, validate *validator.Validate) TaskService {
	return &TaskServiceImpl{
		TaskRepository: taskRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *TaskServiceImpl) Create(ctx context.Context, req web.TaskCreateRequest) (web.TaskResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.TaskResponse{}, err
	}

	if req.Status == nil {
		var pending = domain.StatusPending
		req.Status = &pending
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.TaskResponse{}, err
	}

	task := domain.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
	}

	createdTask, err := service.TaskRepository.Create(ctx, tx, task)
	if err != nil {
		tx.Rollback()
		return web.TaskResponse{}, err
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.TaskResponse{}, errCommit
	}

	return helper.ToTaskResponse(createdTask), nil
}
