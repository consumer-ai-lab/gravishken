package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

func (admin *Admin) GetCollectionName() string {
	return "admins"
}

type AdminRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func FindAdminByUsername(collection *mongo.Collection, username string) (*Admin, error) {
	filter := bson.M{"username": username}

	var admin Admin
	err := collection.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
