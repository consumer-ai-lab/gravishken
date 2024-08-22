package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
)

func BatchRoutes(allControllers *controllers.Class, route *gin.Engine){
	batchRoute := route.Group("/batch")

	batchRoute.POST("/add", func(ctx *gin.Context) {
		var filePath string
		if err := ctx.ShouldBindJSON(&filePath); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		
		allControllers.AddBatchToDB(ctx, filePath)
	})
}