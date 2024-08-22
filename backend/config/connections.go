package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection() (*mongo.Client, error) {
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client, nil
}
