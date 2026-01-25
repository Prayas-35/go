package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Prayas-35/fiber/internal/database"
	"github.com/Prayas-35/fiber/internal/models"
	"github.com/Prayas-35/fiber/utils/helpers"
)

func RegisterUser(user *models.User) error {
	// Hash the password before saving
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.UserCollection == nil {
		return mongo.ErrClientDisconnected
	}

	result, err := database.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	_ = result

	return nil
}
