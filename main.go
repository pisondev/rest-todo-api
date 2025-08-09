package main

import (
	"log"
	"rest-todo-api/app"
	"rest-todo-api/controller"
	"rest-todo-api/exception"
	"rest-todo-api/middleware"
	"rest-todo-api/repository"
	"rest-todo-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	server.Post("/api/register", userController.Register)
	server.Post("/api/login", userController.Login)

	taskRoutes := server.Group("/api/tasks", middleware.AuthMiddleware())

	taskRoutes.Post("", taskController.Create)

	err = server.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
