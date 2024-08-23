package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
	Batch "server/src/models/batch"
)

func BatchRoutes(allControllers *controllers.ControllerClass, route *gin.Engine){
	batchRoute := route.Group("/batch")

	batchRoute.POST("/add", func(ctx *gin.Context) {
		var batchData Batch.Batch
		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		
		allControllers.AddBatchToDB(ctx, &batchData)
	})
}

