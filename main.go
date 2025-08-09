package main

import (
	"log"
	"rest-todo-api/app"
	"rest-todo-api/controller"
	"rest-todo-api/exception"
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

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	server.Post("/api/register", userController.Register)
	server.Post("/api/login", userController.Login)

	err = server.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
