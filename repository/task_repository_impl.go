package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/exception"
	"rest-todo-api/model/domain"
)

type TaskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

func (repository *TaskRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, task domain.Task) (domain.Task, error) {
	SQL := "INSERT INTO tasks(user_id, title, description, status, due_date) VALUES (?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, SQL, task.UserID, task.Title, &task.Description, &task.Status, &task.DueDate)
	if err != nil {
		return domain.Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Task{}, err
	}

	task.ID = int(id)

	SQLSelect := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at FROM tasks WHERE id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, task.ID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (repository *TaskRepositoryImpl) FindTasks(ctx context.Context, tx *sql.Tx, taskFilter TaskFilter) ([]domain.Task, error) {
	SQL := "SELECT id, user_id, title, status, due_date, created_at, updated_at FROM tasks WHERE 1=1"

	var args []any
	if taskFilter.Status != nil && *taskFilter.Status != "" {
		SQL += " AND status = ?"
		args = append(args, taskFilter.Status)
	}

	if taskFilter.DueDate != nil {
		if !taskFilter.DueDate.IsZero() {
			SQL += " AND due_date = ?"
			args = append(args, taskFilter.DueDate)
		}
	}

	SQL += " AND user_id = ? AND deleted_at IS NULL"
	args = append(args, taskFilter.UserID)

	rows, err := tx.QueryContext(ctx, SQL, args...)
	if err != nil {
		return []domain.Task{}, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		task := domain.Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return []domain.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repository *TaskRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, taskID int, userID int) (domain.Task, error) {
	SQL := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at from tasks WHERE id = ? AND user_id = ? AND deleted_at IS NULL"
	rows, err := tx.QueryContext(ctx, SQL, taskID, userID)
	if err != nil {
		return domain.Task{}, err
	}
	defer rows.Close()

	task := domain.Task{}
	if rows.Next() {
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return domain.Task{}, err
		}
		return task, nil
	} else {
		return task, exception.ErrNotFoundTask
	}
}

func (repository *TaskRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, task domain.Task) (domain.Task, error) {
	SQL := "UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ? AND deleted_at IS NULL"
	_, err := tx.ExecContext(ctx, SQL, task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		return domain.Task{}, err
	}

	SQLSelect := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at FROM tasks WHERE id = ? AND deleted_at IS NULL"
	err = tx.QueryRowContext(ctx, SQLSelect, task.ID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (repository *TaskRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, taskID int, userID int) error {
	SQL := "UPDATE tasks SET deleted_at = NOW() WHERE id = ? AND user_id = ? AND deleted_at IS NULL"
	result, err := tx.ExecContext(ctx, SQL, taskID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return exception.ErrNotFoundTask
	}

	return nil
}
