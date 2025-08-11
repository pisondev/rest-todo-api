package service

import (
	"context"
	"database/sql"
	"rest-todo-api/exception"
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
	}

	if req.DueDate != nil {
		if *req.DueDate != "" {
			parsedTime, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return web.TaskResponse{}, exception.ErrBadRequestTimeFormat
			}
			task.DueDate = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
	}

	createdTask, err := service.TaskRepository.Create(ctx, tx, task)
	if err != nil {
		tx.Rollback()
		return web.TaskResponse{}, err
	}
	createdTask.DueDate = sql.NullTime{
		Time:  createdTask.DueDate.Time.UTC().Truncate(time.Second),
		Valid: true,
	}
	createdTask.CreatedAt = createdTask.CreatedAt.UTC().Truncate(time.Second)
	createdTask.UpdatedAt = createdTask.UpdatedAt.UTC().Truncate(time.Second)

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.TaskResponse{}, errCommit
	}

	return helper.ToTaskResponse(createdTask), nil
}

func (service *TaskServiceImpl) FindTasks(ctx context.Context, req web.TaskFilterRequest) ([]web.TaskResponse, error) {
	taskFilter := repository.TaskFilter{}
	tx, err := service.DB.Begin()
	if err != nil {
		return []web.TaskResponse{}, err
	}

	taskFilter.UserID = req.UserID

	// if there is query param, we will include it here
	if req.Status != nil {
		if *req.Status != "" {
			taskFilter.Status = req.Status
		}
	}

	if req.DueDate != nil {
		if *req.DueDate != "" {
			dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return []web.TaskResponse{}, exception.ErrBadRequestTimeFormat
			}

			taskFilter.DueDate = &dueDate
		}
	}

	tasks, err := service.TaskRepository.FindTasks(ctx, tx, taskFilter)
	if err != nil {
		tx.Rollback()
		return []web.TaskResponse{}, err
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return []web.TaskResponse{}, errCommit
	}

	return helper.ToTaskResponses(tasks), nil
}

func (service *TaskServiceImpl) FindByID(ctx context.Context, taskID int, userID int) (web.TaskResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return web.TaskResponse{}, err
	}
	task, err := service.TaskRepository.FindByID(ctx, tx, taskID, userID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.TaskResponse{}, errRollback
		}
		return web.TaskResponse{}, err
	}

	task.DueDate = sql.NullTime{
		Time:  task.DueDate.Time.UTC().Truncate(time.Second),
		Valid: true,
	}
	task.CreatedAt = task.CreatedAt.UTC().Truncate(time.Second)
	task.UpdatedAt = task.UpdatedAt.UTC().Truncate(time.Second)

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.TaskResponse{}, errCommit
	}
	return helper.ToTaskResponse(task), nil
}

func (service *TaskServiceImpl) Update(ctx context.Context, req web.TaskUpdateRequest) (web.TaskResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.TaskResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.TaskResponse{}, err
	}

	selectedTask, err := service.TaskRepository.FindByID(ctx, tx, req.ID, req.UserID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.TaskResponse{}, errRollback
		}
		return web.TaskResponse{}, err
	}

	if selectedTask.UserID != req.UserID {
		return web.TaskResponse{}, exception.ErrUnauthorized
	}

	selectedTask.DueDate = sql.NullTime{
		Time:  selectedTask.DueDate.Time.UTC().Truncate(time.Second),
		Valid: true,
	}
	selectedTask.CreatedAt = selectedTask.CreatedAt.UTC().Truncate(time.Second)
	selectedTask.UpdatedAt = selectedTask.UpdatedAt.UTC().Truncate(time.Second)

	//if even just 1 param was changed, do Update method
	if (req.Title != nil) || (req.Description != nil) || (req.Status != nil) {
		if req.Title != nil {
			selectedTask.Title = *req.Title
		}
		if req.Description != nil {
			selectedTask.Description = req.Description
		}
		if req.Status != nil {
			if (*req.Status != domain.StatusPending) && (*req.Status != domain.StatusDone) {
				return web.TaskResponse{}, exception.ErrBadRequestTaskStatus
			}
			selectedTask.Status = req.Status
		}

		updatedTask, err := service.TaskRepository.Update(ctx, tx, selectedTask)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return web.TaskResponse{}, errRollback
			}
			return web.TaskResponse{}, err
		}

		updatedTask.DueDate = sql.NullTime{
			Time:  updatedTask.DueDate.Time.UTC().Truncate(time.Second),
			Valid: true,
		}
		updatedTask.CreatedAt = updatedTask.CreatedAt.UTC().Truncate(time.Second)
		updatedTask.UpdatedAt = updatedTask.UpdatedAt.UTC().Truncate(time.Second)

		errCommit := tx.Commit()
		if errCommit != nil {
			return web.TaskResponse{}, errCommit
		}
		return helper.ToTaskResponse(updatedTask), nil
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.TaskResponse{}, errCommit
	}
	return helper.ToTaskResponse(selectedTask), nil
}

func (service *TaskServiceImpl) Delete(ctx context.Context, taskID int, userID int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	err = service.TaskRepository.Delete(ctx, tx, taskID, userID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}
