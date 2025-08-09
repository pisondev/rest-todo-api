package service

import (
	"context"
	"database/sql"
	"os"
	"rest-todo-api/exception"
	"rest-todo-api/helper"
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
	"rest-todo-api/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

	_, err = service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != sql.ErrNoRows {
		return web.UserResponse{}, exception.ErrConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.UserResponse{}, err
	}

	user := domain.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now().UTC().Truncate(time.Second),
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

func (service *UserServiceImpl) Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserLoginResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserLoginResponse{}, err
	}
	defer tx.Rollback()

	foundUser, err := service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		return web.UserLoginResponse{}, exception.ErrUnauthorizedLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(req.Password))
	if err != nil {
		return web.UserLoginResponse{}, exception.ErrUnauthorizedLogin
	}

	claims := web.JWTClaims{
		UserID: foundUser.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return web.UserLoginResponse{}, err
	}

	return web.UserLoginResponse{
		Token: tokenString,
	}, nil
}
