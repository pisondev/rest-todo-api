package controller

import (
	"rest-todo-api/model/web"
	"rest-todo-api/service"

	"github.com/gofiber/fiber/v2"
)

type TaskControllerImpl struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return &TaskControllerImpl{
		TaskService: taskService,
	}
}

func (controller *TaskControllerImpl) Create(ctx *fiber.Ctx) error {
	taskCreateRequest := web.TaskCreateRequest{}
	err := ctx.BodyParser(&taskCreateRequest)
	if err != nil {
		return err
	}

	userIDString := ctx.Locals("userID")
	userID, ok := userIDString.(int)
	if !ok {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   "userID not found in the context",
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(webResponse)
	}

	taskCreateRequest.UserID = userID

	taskResponse, err := controller.TaskService.Create(ctx.Context(), taskCreateRequest)
	if err != nil {
		return err
	}

	webResponse := web.WebResponse{
		Code:   201,
		Status: "Created",
		Data:   taskResponse,
	}

	return ctx.Status(fiber.StatusCreated).JSON(webResponse)
}
