package controller

import "github.com/gofiber/fiber/v2"

type TaskController interface {
	Create(ctx *fiber.Ctx) error
	FindTasks(ctx *fiber.Ctx) error
}
