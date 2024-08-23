package controllers

import "go.mongodb.org/mongo-driver/mongo"

type ControllerClass struct {
	Client          *mongo.Client
	AdminCollection *mongo.Collection
	UserCollection  *mongo.Collection
	TestCollection  *mongo.Collection
	BatchCollection *mongo.Collection
}
