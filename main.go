package main

import (
	"log"
	"os"
	"rest-todo-api/app"
	"rest-todo-api/controller"
	"rest-todo-api/exception"
	"rest-todo-api/middleware"
	"rest-todo-api/repository"
	"rest-todo-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	// NEW: Import the CORS middleware for Fiber
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	logger := logrus.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository(logger)
	userService := service.NewUserService(userRepository, db, validate, logger)
	userController := controller.NewUserController(userService, logger)

	taskRepository := repository.NewTaskRepository(logger)
	taskService := service.NewTaskService(taskRepository, db, validate, logger)
	taskController := controller.NewTaskController(taskService, logger)

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	server.Post("/api/register", userController.Register)
	server.Post("/api/login", userController.Login)

	taskRoutes := server.Group("/api/tasks", middleware.AuthMiddleware())

	taskRoutes.Post("", taskController.Create)
	taskRoutes.Get("", taskController.FindTasks)
	taskRoutes.Get("/:taskID", taskController.FindByID)
	taskRoutes.Patch("/:taskID", taskController.Update)
	taskRoutes.Delete("/:taskID", taskController.Delete)

	err = server.Listen(serverPort)
	if err != nil {
		panic(err)
	}
}
