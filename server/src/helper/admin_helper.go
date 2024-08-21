package helper

import (
	"context"
	"fmt"
	"server/src/models/admin"
	"server/src/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(Collection *mongo.Collection, Admin types.ModelInterface) error {

	password := Admin.(*admin.Admin).Password

	fmt.Printf("Original password: %s\n", password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}

	Admin.(*admin.Admin).Password = string(hashedPassword)

	Add_Model_To_DB(Collection, Admin)
	return nil
}

func AdminLogin(Collection *mongo.Collection, Admin types.ModelInterface) error {
	username := Admin.(*admin.Admin).Username
	password := Admin.(*admin.Admin).Password
	secretKey := []byte("token")

	var user admin.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("admin not found")
		}
		return fmt.Errorf("error finding admin: %v", err)
	}

	fmt.Printf("Admin: %v | Password: %s\n", user, password)

	// Compare the hashed password with the plaintext password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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

	Admin.(*admin.Admin).Token = append(Admin.(*admin.Admin).Token, tokenString)

	err = Update_Model_By_ID(Collection, user.ID.Hex(), Admin)

	if err != nil {
		return fmt.Errorf("error updating admin: %v", err)
	}

	fmt.Println("Login successful")

	return nil
}

func AdminLogout(Collection *mongo.Collection, AdminRequest *admin.AdminRequest) error {
	token := AdminRequest.Token
	username := AdminRequest.Username

	var admin admin.Admin
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
