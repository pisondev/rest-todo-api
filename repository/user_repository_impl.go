package repository

import (
	"context"
	"database/sql"
	"rest-todo-api/model/domain"

	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewUserRepository(logger *logrus.Logger) UserRepository {
	return &UserRepositoryImpl{Logger: logger}
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
	repository.Logger.Info("-----START REPO LAYER-----")
	repository.Logger.Info("query sql")
	SQL := "SELECT id, username, hashed_password FROM users WHERE username = ?"
	rows, err := tx.QueryContext(ctx, SQL, username)
	if err != nil {
		repository.Logger.Errorf("failed to query context: %v", err)
		return domain.User{}, err
	}
	defer rows.Close()

	var user domain.User
	repository.Logger.Infof("try to check rows next. Error: %v", err)
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.HashedPassword)
		if err != nil {
			repository.Logger.Errorf("failed to scan row:%v", err)
			return domain.User{}, err
		}
	} else {
		err = sql.ErrNoRows
		repository.Logger.Errorf("no next result row or an error happened while preparing rows.Next: %v", err)
		return domain.User{}, err
	}

	repository.Logger.Infof("final error: %v", err)
	repository.Logger.Info("returned error nil")
	repository.Logger.Info("-----END REPO LAYER-----")
	return user, nil
}
