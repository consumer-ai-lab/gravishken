package types

import (
	"backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type ModelInterface interface {
	GetCollectionName() string
}

var _ ModelInterface = (*models.Admin)(nil)
var _ ModelInterface = (*models.Batch)(nil)
var _ ModelInterface = (*models.User)(nil)
var _ ModelInterface = (*models.UserTest)(nil)
var _ ModelInterface = (*models.Test)(nil)

// Define your JWT claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
