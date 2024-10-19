package route

import (
	"common/models/user"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"net/http"
	"regexp"
	"server/src/controllers"
	"strconv"
	"strings"

	middleware "server/src/middleware"
)

func UserRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	userRoute := route.Group("/user")

	userRoute.POST("/login", func(ctx *gin.Context) {
		var userModel user.UserLoginRequest
		if err := ctx.ShouldBindJSON(&userModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UserLoginHandler(ctx, &userModel)
	})

	authenticated := userRoute.Group("/")

	authenticated.Use(middleware.AdminJWTAuthMiddleware(allControllers.UserCollection))

	authenticated.GET("/get_all_users", func(ctx *gin.Context) {
		var users []user.User
		cursor, err := allControllers.UserCollection.Find(context.Background(), bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &users); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"users": users})
	})

	authenticated.GET("/paginated_users", func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
		search := strings.TrimSpace(ctx.Query("search"))

		fmt.Printf("Received request - Page: %d, Limit: %d, Search: %s\n", page, limit, search)

		filter := bson.M{}
		if search != "" {
			// Escape special regex characters and use case-insensitive search
			escapedSearch := regexp.QuoteMeta(search)
			filter = bson.M{
				"$or": []bson.M{
					{"username": primitive.Regex{Pattern: escapedSearch, Options: "i"}}, // Changed from name to username
					{"batch": primitive.Regex{Pattern: escapedSearch, Options: "i"}},    // Changed from batchName to batch
				},
			}
		}

		totalUsers, err := allControllers.UserCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			fmt.Printf("Error counting users: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total users"})
			return
		}

		totalPages := int(math.Ceil(float64(totalUsers) / float64(limit)))

		if page < 1 {
			page = 1
		} else if page > totalPages && totalPages > 0 {
			page = totalPages
		}

		skip := (page - 1) * limit

		opts := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "username", Value: 1}}) // Changed from name to username

		var users []user.User
		cursor, err := allControllers.UserCollection.Find(context.Background(), filter, opts)
		if err != nil {
			fmt.Printf("Error fetching users: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &users); err != nil {
			fmt.Printf("Error decoding users: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users":       users,
			"totalPages":  totalPages,
			"currentPage": page,
			"totalUsers":  totalUsers,
		})
	})

	userRoute.PUT("/update_user", func(ctx *gin.Context) {
		var updateRequest user.UserModelUpdateRequest
		if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Default().Println("User Request: ", updateRequest)

		err := allControllers.UpdateUser(ctx, &updateRequest)

		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in updating user!",
				"error":   err,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "user updated successfully!!",
		})
	})

	userRoute.DELETE("/delete_user", func(ctx *gin.Context) {
		var deleteRequest struct {
			UserId string `json:"userId"`
		}
		if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if deleteRequest.UserId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		err := allControllers.DeleteUser(ctx, deleteRequest.UserId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

}
