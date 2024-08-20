package helper

import (
	"backend/auth"
	"backend/models"
	"backend/types"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(Collection *mongo.Collection, Admin types.ModelInterface) error {

	password := Admin.(*models.Admin).Password

	fmt.Printf("Original password: %s\n", password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}

	Admin.(*models.Admin).Password = string(hashedPassword)

	Add_Model_To_DB(Collection, Admin)
	return nil
}

func AdminLogin(Collection *mongo.Collection, Admin types.ModelInterface) error {
	username := Admin.(*models.Admin).Username
	password := Admin.(*models.Admin).Password
	secretKey := []byte("token")

	var admin models.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("admin not found")
		}
		return fmt.Errorf("error finding admin: %v", err)
	}

	fmt.Printf("Admin: %v | Password: %s\n", admin, password)

	// Compare the hashed password with the plaintext password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %v", err)
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
		return err
	}

	Admin.(*models.Admin).Token = append(Admin.(*models.Admin).Token, tokenString)

	err = Update_Model_By_ID(Collection, admin.ID.Hex(), Admin)

	if err != nil {
		return fmt.Errorf("error updating admin: %v", err)
	}

	fmt.Println("Login successful")

	return nil
}

func AdminLogout(Collection *mongo.Collection, Admin *models.AdminRequest) error {
	token := Admin.Token
	username := Admin.Username

	var admin models.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
	if err != nil {
		return err
	}

	new_token := []string{}

	for _, t := range admin.Token {
		if t != token {
			new_token = append(new_token, t)
		}
	}

	admin.Token = new_token
	result, err := Collection.ReplaceOne(context.TODO(), bson.M{"username": username}, admin)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID: %s", username)
	}

	fmt.Println("Admin Logout successfully")

	return nil
}

func AuthenticateAdmin(Collection *mongo.Collection, Admin *models.AdminRequest) bool {
	token := Admin.Token

	verified, err := auth.TokenVerifier(Collection, token)
	fmt.Println(verified)

	return err == nil

}

func ChangePassword(Collection *mongo.Collection, model *models.AdminChangePassword) error {

	username := model.Username

	var ADMIN models.Admin
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
