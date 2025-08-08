package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO users(username, hashed_password, created_at) VALUES (?,?,?)"
	result, err := tx.ExecContext(ctx, SQL, user.Username, user.HashedPassword, user.CreatedAt)
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

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	SQL := "SELECT id, username, hashed_password FROM users WHERE username = ?"
	rows, err := tx.QueryContext(ctx, SQL, username)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	var user domain.User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.HashedPassword)
		if err != nil {
			return domain.User{}, err
		}
	}
	return user, nil
}
