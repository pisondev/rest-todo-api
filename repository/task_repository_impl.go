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
	SQL := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at from tasks WHERE id = ? AND user_id = ?"
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
	SQL := "UPDATE tasks SET title = ?, description = ?, status = ?, updated_at = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	if err != nil {
		return domain.Task{}, err
	}

	// SQL := "UPDATE tasks SET id = id"
	// var args []any
	// if task.Title != "" {
	// 	SQL += ", title = ?"
	// 	args = append(args, task.Title)
	// }
	// if task.Description != nil {
	// 	if *task.Description != "" {
	// 		SQL += ", description = ?"
	// 		args = append(args, task.Description)
	// 	}
	// }
	// if task.Status != nil {
	// 	if *task.Status != "" {
	// 		SQL += ", status = ?"
	// 		args = append(args, task.Status)
	// 	}
	// }
	// SQL += " WHERE id = ?"
	// args = append(args, task.ID)

	// _, err := tx.QueryContext(ctx, SQL, args...)
	// if err != nil {
	// 	return domain.Task{}, err
	// }

	return task, nil
}
