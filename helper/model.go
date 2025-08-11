package helper

import (
	"database/sql"
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
	"time"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func ToTaskResponse(task domain.Task) web.TaskResponse {
	response := web.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.DueDate.Valid {
		response.DueDate = &task.DueDate.Time
	}
	return response
}

func ToTaskResponses(tasks []domain.Task) []web.TaskResponse {
	var taskResponses []web.TaskResponse
	for _, task := range tasks {
		task.DueDate = sql.NullTime{
			Time:  task.DueDate.Time.UTC().Truncate(time.Second),
			Valid: true,
		}
		task.CreatedAt = task.CreatedAt.UTC().Truncate(time.Second)
		task.UpdatedAt = task.UpdatedAt.UTC().Truncate(time.Second)

		taskResponses = append(taskResponses, ToTaskResponse(task))
	}
	return taskResponses
}
