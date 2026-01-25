package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	ServerPort string
	JWTSecret  string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080" // default port if none is set
		}
	}

	return Config{
		MongoURI:   os.Getenv("MONGO_URI"),
		ServerPort: port,
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
