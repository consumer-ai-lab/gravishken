package types

import (
	"common/models/admin"
	"common/models/batch"
	"common/models/test"
	"common/models/user"

	"github.com/golang-jwt/jwt/v5"
)

type ModelInterface interface {
	GetCollectionName() string
}

var _ ModelInterface = (*admin.Admin)(nil)
var _ ModelInterface = (*batch.Batch)(nil)
var _ ModelInterface = (*user.User)(nil)
var _ ModelInterface = (*user.UserSubmission)(nil)
var _ ModelInterface = (*test.BatchTests)(nil)

// Define your JWT claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
