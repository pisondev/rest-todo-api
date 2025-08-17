package controller

import (
	"rest-todo-api/model/domain"
	"rest-todo-api/model/web"
	"rest-todo-api/service"
	"strconv"

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
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
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

	return ctx.Status(fiber.StatusCreated).JSON(taskResponse)
}

func (controller *TaskControllerImpl) FindTasks(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	dueDateStr := ctx.Query("due_date")

	userIDStr := ctx.Locals("userID")
	userID, ok := userIDStr.(int)
	if !ok {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   "userID not found in the context",
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(webResponse)
	}

	taskFilterRequest := web.TaskFilterRequest{
		Status:  (*domain.TaskStatus)(&status),
		UserID:  userID,
		DueDate: &dueDateStr,
	}
	taskResponses, err := controller.TaskService.FindTasks(ctx.Context(), taskFilterRequest)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(taskResponses)
}

func (controller *TaskControllerImpl) FindByID(ctx *fiber.Ctx) error {
	taskIDStr := ctx.Params("taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return err
	}

	userIDString := ctx.Locals("userID")
	userID := userIDString.(int)

	taskResponse, err := controller.TaskService.FindByID(ctx.Context(), taskID, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(taskResponse)
}

func (controller *TaskControllerImpl) Update(ctx *fiber.Ctx) error {
	taskUpdateRequest := web.TaskUpdateRequest{}
	err := ctx.BodyParser(&taskUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(webResponse)
	}

	taskIDStr := ctx.Params("taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return err
	}
	taskUpdateRequest.ID = taskID

	userID := ctx.Locals("userID")
	taskUpdateRequest.UserID = userID.(int)

	taskResponse, err := controller.TaskService.Update(ctx.Context(), taskUpdateRequest)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(taskResponse)
}

func (controller *TaskControllerImpl) Delete(ctx *fiber.Ctx) error {
	taskIDStr := ctx.Params("taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return err
	}
	userIDStr := ctx.Locals("userID")
	userID, ok := userIDStr.(int)
	if !ok {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   "userID not found in the context",
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(webResponse)
	}

	err = controller.TaskService.Delete(ctx.Context(), taskID, userID)
	if err != nil {
		return err
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}
