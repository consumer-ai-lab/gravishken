package route

import (
	"server/src/controllers"
	User "server/src/models/user"
	Batch "server/src/models/batch"

	"github.com/gin-gonic/gin"
)

func UserRoutes(allControllers *controllers.ControllerClass, route *gin.Engine){
	userRoute := route.Group("/user")

	userRoute.POST("/login", func(ctx *gin.Context) {
		var userModel User.UserLoginRequest
		if err := ctx.ShouldBindJSON(&userModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UserLoginHandler(ctx, &userModel)
	})

	userRoute.POST("/create_batch_data", func(ctx *gin.Context) {
		var batchData Batch.Batch
		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

	})

	userRoute.POST("/update_user_data", func(ctx *gin.Context) {
		var userUpdateRequest User.UserUpdateRequest

		if err := ctx.ShouldBindJSON(&userUpdateRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UpdateUserData(ctx, &userUpdateRequest)

	})
}


