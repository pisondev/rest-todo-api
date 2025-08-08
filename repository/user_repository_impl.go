package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/model/domain"
)

type UserRepositoryImpl struct {
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO users(username, hashed_password) VALUES (?,?)"
	result, err := tx.ExecContext(ctx, SQL, user.Username, user.HashedPassword)
	if err != nil {
		return domain.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}

	user.Id = int(id)
	return user, nil
}
