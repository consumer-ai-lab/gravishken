package route

import (
	User "common/models/user"
	"server/src/controllers"

	"github.com/gin-gonic/gin"
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
}
