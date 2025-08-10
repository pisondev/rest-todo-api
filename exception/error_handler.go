package exception

import (
	"errors"
	"rest-todo-api/model/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var validationErr validator.ValidationErrors

	code := fiber.StatusInternalServerError
	status := "INTERNAL SERVER ERROR"

	if errors.As(err, &validationErr) {
		code = fiber.StatusBadRequest
		status = "BAD REQUEST"
	}
	if errors.Is(err, ErrBadRequestTimeFormat) {
		code = fiber.StatusBadRequest
		status = "BAD REQUEST"
	}
	if errors.Is(err, ErrUnauthorized) {
		code = fiber.StatusUnauthorized
		status = "UNAUTHORIZED"
	}
	if errors.Is(err, ErrUnauthorizedLogin) {
		code = fiber.StatusUnauthorized
		status = "UNAUTHORIZED"
	}
	if errors.Is(err, ErrConflict) {
		code = fiber.StatusConflict
		status = "CONFLICT"
	}

	webResponse := web.WebResponse{
		Code:   code,
		Status: status,
		Data:   err.Error(),
	}

	return ctx.Status(code).JSON(webResponse)
}
