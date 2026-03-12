package main

import (
	"log"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/controllers"
	_ "github.com/LanangDepok/project-management/docs"
	"github.com/LanangDepok/project-management/repositories"
	"github.com/LanangDepok/project-management/routes"
	"github.com/LanangDepok/project-management/services"
	"github.com/gofiber/fiber/v3"
)

// @title           Project Management API
// @version         1.0
// @description     REST API for Project Management built with Go Fiber v3 and GORM.
// @host            localhost:8000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and JWT token.
func main() {
	config.LoadEnv()
	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName: "Project Management API v1.0",
	})

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	boardRepo := repositories.NewBoardRepository()
	boardService := services.NewBoardService(boardRepo, userRepo)
	boardController := controllers.NewBoardController(boardService)

	routes.Setup(app, userController, boardController)

	port := config.AppConfig.AppPort
	log.Printf("Server running on http://localhost:%s", port)
	log.Printf("Swagger UI at http://localhost:%s/swagger/", port)
	log.Fatal(app.Listen(":" + port))
}
