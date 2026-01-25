package routes

import (
	"github.com/Prayas-35/fiber/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {
	authGroup := app.Group("/auth")

	authGroup.Post("/register", controllers.RegisterUserController)
}
