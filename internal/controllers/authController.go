package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Prayas-35/fiber/internal/middlewares"
	"github.com/Prayas-35/fiber/internal/models"
	"github.com/Prayas-35/fiber/internal/services"
	"github.com/Prayas-35/fiber/utils/helpers"
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

	token, err := helpers.GenerateJWT("", user.ID.Hex(), user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID.Hex(), "username": user.Username, "email": user.Email, "token": token})
}

func LogInUserController(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password are required"})
	}

	user, token, err := services.LogInUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
		}
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "database unavailable"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to log in"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": user.ID.Hex(), "username": user.Username, "email": user.Email, "token": token})
}

func GetUserProfileController(c *fiber.Ctx) error {
	authUser, ok := c.Locals("authUser").(middlewares.AuthUser)
	if !ok || authUser.ID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	user, err := services.GetUserProfile(authUser.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "database unavailable"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve user profile"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": user.ID.Hex(), "username": user.Username, "email": user.Email, "created_at": user.CreatedAt, "updated_at": user.UpdatedAt})
}
