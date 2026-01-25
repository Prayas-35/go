package routes

import (
	"github.com/Prayas-35/fiber/internal/controllers"
	"github.com/Prayas-35/fiber/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router, jwtSecret string) {
	authGroup := app.Group("/auth")

	authGroup.Post("/register", controllers.RegisterUserController)
	authGroup.Post("/login", controllers.LogInUserController)
	authGroup.Get("/profile", middlewares.VerifyTokenMiddleware(), controllers.GetUserProfileController)
}
