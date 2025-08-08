package service

import (
	"context"
	"database/sql"
	"rest-todo-api/helper"
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
	"rest-todo-api/repository"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req web.UserAuthRequest) (web.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.UserResponse{}, err
	}

	user := domain.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
	}

	savedUser, err := service.UserRepository.Save(ctx, tx, user)
	if err != nil {
		return web.UserResponse{}, err
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.UserResponse{}, errCommit
	}

	return helper.ToUserResponse(savedUser), nil
}
