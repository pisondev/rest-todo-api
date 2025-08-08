package main

import (
	"rest-todo-api/app"
	"rest-todo-api/controller"
	"rest-todo-api/exception"
	"rest-todo-api/repository"
	"rest-todo-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	server.Post("/api/register", userController.Register)
	server.Post("/api/login", userController.Login)

	err := server.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
