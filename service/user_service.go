package service

import (
	"context"
	"rest-todo-api/model/web"
)

type UserService interface {
	Register(ctx context.Context, req web.UserAuthRequest) (web.UserResponse, error)
	Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error)
}
