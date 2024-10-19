package models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return client.Database(dbName).Collection(collectionName)
}
