package helper

import (
	"common/models/admin"
	"context"
	"fmt"
	"server/src/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(Collection *mongo.Collection, Admin types.ModelInterface) error {

	password := Admin.(*admin.Admin).Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}

	Admin.(*admin.Admin).Password = string(hashedPassword)

	Add_Model_To_DB(Collection, Admin)
	return nil
}

func AdminLogin(Collection *mongo.Collection, Admin types.ModelInterface) (string, error) {
	username := Admin.(*admin.Admin).Username
	password := Admin.(*admin.Admin).Password
	secretKey := []byte("token")

	var user admin.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error finding admin: %v", err)
	}

	// Compare the hashed password with the plaintext password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(48 * time.Hour)

	claims := &types.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing the token:", err)
		return "", err
	}

	fmt.Println("Login successful")

	return tokenString, nil
}

func ChangePassword(Collection *mongo.Collection, model *admin.AdminChangePassword) error {

	username := model.Username

	var ADMIN admin.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&ADMIN)

	if err != nil {
		return err
	}

	password := ADMIN.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	ADMIN.Password = string(hashedPassword)

	_, err = Collection.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": ADMIN})

	if err != nil {
		return err
	}

	fmt.Println("Password changed successfully")
	return nil

}
