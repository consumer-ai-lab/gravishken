package route

import (
	"common/models/user"
	"context"
	"log"
	"math"
	"net/http"
	"server/src/controllers"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

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
		skip := (page - 1) * limit

		var users []user.User

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
		cursor, err := allControllers.UserCollection.Find(context.Background(), bson.M{}, opts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &users); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		totalUsers, err := allControllers.UserCollection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total users"})
			return
		}

		totalPages := int(math.Ceil(float64(totalUsers) / float64(limit)))

		ctx.JSON(http.StatusOK, gin.H{
			"users":       users,
			"totalPages":  totalPages,
			"currentPage": page,
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
