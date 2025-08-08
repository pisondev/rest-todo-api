package helper

import (
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
