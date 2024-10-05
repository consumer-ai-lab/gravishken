package route

import (
	"common/models/user"
	"github.com/gin-gonic/gin"
	"server/src/controllers"
	"net/http"
	"context"
	"strconv"
	"math"
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

	
}
