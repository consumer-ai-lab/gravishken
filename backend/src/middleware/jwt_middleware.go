package middleware

import (
	"log"
	"net/http"
	"server/src/auth"
	"server/src/helper"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		claims, err := auth.ApplicationTokenVerifier(userCollection, tokenString)
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
