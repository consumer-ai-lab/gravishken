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
	secretKey := []byte("TODO:add-a-secret-key-from-env") 

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
		return "", fmt.Errorf("invalid credentials")
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
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return tokenString, nil
}



func UpdateTypingTestText(collection *mongo.Collection, testID string, typingText string) error {
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": testID, "type": "typing"},
		bson.M{"$set": bson.M{"typingText": typingText}},
	)
	if err != nil {
		return fmt.Errorf("error updating typing test text: %v", err)
	}

	return nil
}
