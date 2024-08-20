package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Username string             `bson:"username" json:"username"`
    Password string             `bson:"password" json:"-"`
    Token    []string           `bson:"token" json:"token"`
}

func (admin *Admin) GetCollectionName() string {
    return "admins"
}


type AdminRequest struct {
    Username string `json:"username"`
    Token string `json:"token"`
}

type AdminChangePassword struct {
    Username string `json:"username"`
    NewPassword string `json:"oldPassword"`
}

