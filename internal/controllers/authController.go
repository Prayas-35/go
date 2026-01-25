package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Prayas-35/fiber/internal/models"
	"github.com/Prayas-35/fiber/internal/services"
)

func RegisterUserController(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username, email, and password are required"})
	}

	now := primitive.NewDateTimeFromTime(time.Now().UTC())
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := services.RegisterUser(&user); err != nil {
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "database unavailable"})
		}

		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user already exists"})
				}
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID.Hex()})
}
