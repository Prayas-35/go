package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UserIndexes() {
	if UserCollection == nil {
		return
	}

	// Create unique index on email field in UserCollection
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := UserCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{emailIndexModel})
	if err != nil {
		log.Fatal("Failed to create index on email field:", err)
	}
}

func InitIndexes() {
	UserIndexes()
}
