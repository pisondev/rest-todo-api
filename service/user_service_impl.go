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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Logger         *logrus.Logger
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate, logger *logrus.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		Logger:         logger,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req web.UserAuthRequest) (web.UserResponse, error) {
	service.Logger.Info("-----START SERVICE LAYER-----")
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	service.Logger.Infof("calling FindByUsername. Error: %v", err)
	_, err = service.UserRepository.FindByUsername(ctx, tx, req.Username)
	service.Logger.Infof("FindByUsername has done. Error: %v", err)
	if err != sql.ErrNoRows {
		service.Logger.Error("username already exist")
		service.Logger.Info("rollback tx")
		errRollback := tx.Rollback()
		if errRollback != nil {
			service.Logger.Errorf("failed to rollback tx: %v", errRollback)
			return web.UserResponse{}, errRollback
		}
		return web.UserResponse{}, exception.ErrConflict
	}

	service.Logger.Info("generate hash from password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		service.Logger.Errorf("failed to generate from password: %v", err)
		return web.UserResponse{}, err
	}

	user := domain.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now().UTC().Truncate(time.Second),
	}

	service.Logger.Info("call Save Repository")
	savedUser, err := service.UserRepository.Save(ctx, tx, user)
	if err != nil {
		service.Logger.Errorf("failed to use Save Repository: %v", err)
		service.Logger.Info("rollback tx")
		errRollback := tx.Rollback()
		if errRollback != nil {
			service.Logger.Errorf("failed to rollback tx: %v", errRollback)
			return web.UserResponse{}, errRollback
		}

		return web.UserResponse{}, err
	}

	service.Logger.Info("commit tx")
	errCommit := tx.Commit()
	if errCommit != nil {
		service.Logger.Errorf("failed to commit tx: %v", errCommit)
		return web.UserResponse{}, errCommit
	}

	service.Logger.Info("return nil for error")
	service.Logger.Info("-----END SERVICE LAYER-----")
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

	foundUser, err := service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		tx.Rollback()
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

	errCommit := tx.Commit()
	if errCommit != nil {
		return web.UserLoginResponse{}, errCommit
	}

	return web.UserLoginResponse{
		Token: tokenString,
	}, nil
}
