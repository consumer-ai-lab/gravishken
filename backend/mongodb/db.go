package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() ( *mongo.Client , error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Default().Print("Loaded .env file")

	MONGODB_URI := os.Getenv("MONGODB_URI")
	uri := MONGODB_URI
	
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client, nil
}
