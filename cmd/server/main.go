package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Prayas-35/fiber/config"
	"github.com/Prayas-35/fiber/internal/database"
	"github.com/Prayas-35/fiber/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(
		fiber.Config{
			AppName: "FinZ",
		},
	)

	app.Use(logger.New())
	app.Use(etag.New())
	app.Use(cors.New())
	app.Use(cache.New())

	// Load config (e.g., from .env)
	cfg := config.LoadConfig()

	// Connect to MongoDB
	dbClient := database.ConnectMongo(database.Config{
		MongoURI: cfg.MongoURI,
	})
	if dbClient == nil {
		log.Fatal("MongoDB connection failed")
	}
	database.InitCollections()
	database.InitIndexes()
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}
	}()

	// Or extend your config for customization
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		ExposeHeaders: "Content-Length",
		MaxAge:        300,
	}))

	api := app.Group("/api")

	routes.AuthRoutes(api, cfg.JWTSecret)

	app.Listen(":" + cfg.ServerPort)
	log.Println("Server is running on port", cfg.ServerPort)
}
