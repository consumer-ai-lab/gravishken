package route

import (
	"server/src/controllers"
	"server/src/middleware"

	Batch "common/models/batch"

	"github.com/gin-gonic/gin"
)

func BatchRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	batchRoute := route.Group("/batch")
	batchRoute.Use(middleware.UserJWTAuthMiddleware(allControllers.UserCollection))

	batchRoute.POST("/add", func(ctx *gin.Context) {
		var batchData Batch.Batch
		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AddBatchToDB(ctx, &batchData)
	})

	batchRoute.GET("/get_batches", func(ctx *gin.Context) {
		allControllers.GetBatches(ctx)
	})
}
