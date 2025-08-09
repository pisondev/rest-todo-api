package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/model/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, tx *sql.Tx, task domain.Task) (domain.Task, error)
}
