package types

import (
	"server/src/models/admin"
	"server/src/models/batch"
	"server/src/models/test"
	"server/src/models/user"

	"github.com/golang-jwt/jwt/v5"
)

type ModelInterface interface {
	GetCollectionName() string
}

var _ ModelInterface = (*admin.Admin)(nil)
var _ ModelInterface = (*batch.Batch)(nil)
var _ ModelInterface = (*user.User)(nil)
var _ ModelInterface = (*user.UserTest)(nil)
var _ ModelInterface = (*test.Test)(nil)

// Define your JWT claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}