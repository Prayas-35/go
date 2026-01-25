package database

import "go.mongodb.org/mongo-driver/mongo"

var UserCollection *mongo.Collection

func InitCollections() {
	if Client == nil {
		return
	}

	UserCollection = Client.Database("fiber_db").Collection("users")
}
