package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func LogInUser(email, password string) (*models.User, string, error) {
	user := &models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.UserCollection == nil {
		return nil, "", mongo.ErrClientDisconnected
	}

	err := database.UserCollection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(user)
	if err != nil {
		return nil, "", err
	}

	// Verify the password
	if !helpers.CheckPasswordHash(password, user.Password) {
		return nil, "", mongo.ErrNoDocuments
	}

	token, err := helpers.GenerateJWT("", user.ID.Hex(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func GetUserProfile(userID string) (*models.User, error) {
	user := &models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if database.UserCollection == nil {
		return nil, mongo.ErrClientDisconnected
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = database.UserCollection.FindOne(ctx, map[string]interface{}{"_id": objID}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
