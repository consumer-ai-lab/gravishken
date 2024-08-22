package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitAuthRoutes(db *mongo.Client, route *gin.Engine) {
	allControllers := controllers.Class{Client: db}

	groupRoute := route.Group("/auth")
	groupRoute.POST("/login", SampleHandler)
	groupRoute.POST("/register", SampleHandler)
	groupRoute.POST("/refresh", SampleHandler)

	adminRoute := groupRoute.Group("/admin")
	adminRoute.POST("/login", allControllers.AdminLoginHandler)
	adminRoute.POST("/register", SampleHandler)
	adminRoute.POST("/refresh", SampleHandler)
}

func SampleHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Auth routes here",
	})
}
