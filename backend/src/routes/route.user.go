package route

import (
	User "common/models/user"
	"github.com/gin-gonic/gin"
	"server/src/controllers"
	middleware "server/src/middleware"
)

func UserRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	userRoute := route.Group("/user")

	userRoute.POST("/login", func(ctx *gin.Context) {
		var userModel User.UserLoginRequest
		if err := ctx.ShouldBindJSON(&userModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UserLoginHandler(ctx, &userModel)
	})

	authenticated := userRoute.Group("/")
	authenticated.Use(middleware.UserJWTAuthMiddleware(allControllers.UserCollection))

	authenticated.POST("/some_protected_route", func(ctx *gin.Context) {
	})
}
