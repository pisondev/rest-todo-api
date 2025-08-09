package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/model/domain"
)

type TaskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

func (repository *TaskRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, task domain.Task) (domain.Task, error) {
	SQL := "INSERT INTO tasks(user_id, title, description, status, due_date, created_at, updated_at) VALUES (?,?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, SQL, task.UserID, task.Title, &task.Description, &task.Status, &task.DueDate, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return domain.Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, err
	}

	task.ID = int(id)
	return task, nil
}
