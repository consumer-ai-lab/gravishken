package auth

import (
	// "common/models/admin"
	"common/models/user"
	"errors"
	"fmt"
	"os"
	"github.com/golang-jwt/jwt/v5"
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

// TokenVerifier verifies the validity of a token and checks if it is valid for a specific admin.
// It takes the MongoDB collection for admin and a token string as input parameters.
// func TokenVerifier(Collection *mongo.Collection, tokenString string) (jwt.MapClaims, error) {

// 	_, claims, err := VerifyToken(tokenString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	userName := claims["userName"].(string)
// 	operator, err := admin.FindByUsername(Collection, userName)
// 	if err != nil || operator == nil {
// 		return nil, errors.New("operator not found")
// 	}

// 	for _, t := range operator.Token {
// 		if t == tokenString {
// 			return claims, nil
// 		}
// 	}

// 	return nil, errors.New("token not valid for user")
// }

func ApplicationTokenVerifier(Collection *mongo.Collection, tokenString string) (jwt.MapClaims, error) {
	_, claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	userName := claims["username"].(string)
	user, err := user.FindByUsername(Collection, userName)
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

// func AuthenticateAdmin(Collection *mongo.Collection, Admin *admin.AdminRequest) bool {
// 	token := Admin.Token

// 	verified, err := TokenVerifier(Collection, token)
// 	fmt.Println(verified)

// 	return err == nil

// }
