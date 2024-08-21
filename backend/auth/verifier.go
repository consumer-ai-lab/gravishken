package auth

import (
	"gravtest/models"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tokenKey = []byte("token")

func VerifyToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return tokenKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}

func FindAdminByUsername(Collection *mongo.Collection, userName string) (*models.Admin, error) {

	filter := bson.M{"userName": userName}

	var admin models.Admin
	err := Collection.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func FindUserByUsername(Collection *mongo.Collection, userName string) (*models.User, error) {

	filter := bson.M{"username": userName}

	var user models.User
	err := Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func TokenVerifier(Collection *mongo.Collection, tokenString string) (jwt.MapClaims, error) {
	_, claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	userName := claims["userName"].(string)
	operator, err := FindAdminByUsername(Collection, userName)
	if err != nil || operator == nil {
		return nil, errors.New("operator not found")
	}

	for _, t := range operator.Token {
		if t == tokenString {
			return claims, nil
		}
	}

	return nil, errors.New("token not valid for user")
}

func ApplicationTokenVerifier(Collection *mongo.Collection, tokenString string) (jwt.MapClaims, error) {
	_, claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	userName := claims["username"].(string)
	user, err := FindUserByUsername(Collection, userName)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}

func ApiKeyVerifier(apiKey string) (bool, error) {
	backendAPISecret := os.Getenv("BACKEND_API_SECRET")
	if backendAPISecret == "" {
		return false, errors.New("backend API secret is not set")
	}

	fmt.Println("apikey = ", apiKey)
	fmt.Println("api_sec = ", backendAPISecret)

	return apiKey == backendAPISecret, nil
}

func ValidRequestVerifier(Collection *mongo.Collection, tokenString, apiKey string) (bool, error) {
	fmt.Println("validRequestVerifier: called")

	claims, err := ApplicationTokenVerifier(Collection, tokenString)
	if err != nil {
		fmt.Println("validRequestVerifier: decoded error: ", err)
		return false, err
	}
	fmt.Println("validRequestVerifier: decoded: ", claims)

	apiKeyResult, err := ApiKeyVerifier(apiKey)
	if err != nil {
		fmt.Println("validRequestVerifier: apiKeyResult error: ", err)
		return false, err
	}
	fmt.Println("validRequestVerifier: apiKeyResult: ", apiKeyResult)

	return claims != nil && apiKeyResult, nil
}
