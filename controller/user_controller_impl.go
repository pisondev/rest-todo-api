package controller

import (
	"rest-todo-api/helper"
	"rest-todo-api/model/web"
	"rest-todo-api/service"

	"github.com/gofiber/fiber/v2"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Register(ctx *fiber.Ctx) error {
	userAuthRequest := web.UserAuthRequest{}
	err := ctx.BodyParser(&userAuthRequest)
	helper.PanicIfError(err)

	userResponse, err := controller.UserService.Register(ctx.Context(), userAuthRequest)
	helper.PanicIfError(err)
	webResponse := web.WebResponse{
		Code:   201,
		Status: "Created",
		Data:   userResponse,
	}

	return ctx.Status(fiber.StatusCreated).JSON(webResponse)
}
