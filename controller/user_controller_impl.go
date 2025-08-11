package controller

import (
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
	if err != nil {
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	userResponse, err := controller.UserService.Register(ctx.Context(), userAuthRequest)
	if err != nil {
		return err
	}
	webResponse := web.WebResponse{
		Code:   201,
		Status: "Created",
		Data:   userResponse,
	}

	return ctx.Status(fiber.StatusCreated).JSON(webResponse)
}

func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	userAuthRequest := web.UserAuthRequest{}
	err := ctx.BodyParser(&userAuthRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	userResponse, err := controller.UserService.Login(ctx.Context(), userAuthRequest)
	if err != nil {
		return err
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}
