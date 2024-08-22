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
	Token    []string           `bson:"token" json:"token"`
}

func (admin *Admin) GetCollectionName() string {
	return "admins"
}

type AdminRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type AdminChangePassword struct {
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

func FindByUsername(Collection *mongo.Collection, userName string) (*Admin, error) {

	filter := bson.M{"userName": userName}

	var admin Admin
	err := Collection.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
