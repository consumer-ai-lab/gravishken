package main

import (
	"common/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/src/helper"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

var tokenKey = []byte("token")

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return tokenKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token has expired")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ApplicationTokenVerifier(Collection *mongo.Collection, tokenString string) (jwt.MapClaims, error) {
	claims, err := VerifyJWT(tokenString)
	if err != nil {
		return nil, err
	}

	userName := claims["username"].(string)
	user, err := models.FindByUsername(Collection, userName)
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

func UserJWTAuthMiddleware(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := ApplicationTokenVerifier(userCollection, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func AdminJWTAuthMiddleware(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Entering AdminJWTAuthMiddleware")

		token, err := c.Cookie("auth_token")
		log.Printf("Auth token from cookie: %s", token)

		if err != nil {
			log.Printf("Error getting auth_token cookie: %v", err)
			c.JSON(401, gin.H{
				"isAuthenticated": false,
				"error":           "No token found",
			})
			return
		}
		log.Println("Auth token cookie found")

		claims, err := helper.ValidateAdminToken(token)
		if err != nil {
			log.Printf("Admin token validation failed: %v", err)
			c.JSON(401, gin.H{
				"isAuthenticated": false,
				"error":           "Invalid token",
			})
			return
		}
		log.Println("Admin token validated successfully")

		c.Set("claims", claims)
		log.Println("Claims set in context")

		c.Next()
		log.Println("Exiting AdminJWTAuthMiddleware")
	}
}
