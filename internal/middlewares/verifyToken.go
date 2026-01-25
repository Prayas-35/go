package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUser struct {
	ID    string
	Email string
}

func VerifyTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing or malformed token"})
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid token claims"})
		}

		userID, ok := claims["userId"].(string)
		email, ok2 := claims["email"].(string)

		if !ok || !ok2 {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid auth payload"})
		}

		c.Locals("authUser", AuthUser{
			ID:    userID,
			Email: email,
		})

		return c.Next()
	}
}
