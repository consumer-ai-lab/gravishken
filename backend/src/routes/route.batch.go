package route

import (
	"server/src/controllers"
	// "server/src/middleware"

	"common/models/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

func BatchRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	batchRoute := route.Group("/batch")
	// batchRoute.Use(middleware.UserJWTAuthMiddleware(allControllers.UserCollection))

	batchRoute.POST("/add", func(ctx *gin.Context) {
		var batchData struct {
			BatchName     string   `json:"batchName"`
			SelectedTests []string `json:"selectedTests"`
		}

		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		// Convert string IDs to ObjectIDs
		testObjectIDs := make([]primitive.ObjectID, 0, len(batchData.SelectedTests))
		for _, testID := range batchData.SelectedTests {
			objectID, err := primitive.ObjectIDFromHex(testID)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid test ID format"})
				return
			}
			testObjectIDs = append(testObjectIDs, objectID)
		}

		newBatch := batch.Batch{
			Name:  batchData.BatchName,
			Tests: testObjectIDs,
		}

		allControllers.AddBatchToDB(ctx, &newBatch)
		
		ctx.JSON(200, gin.H{"message": "Batch added successfully", "batch": newBatch})
	})

	batchRoute.GET("/get_batches", func(ctx *gin.Context) {
		allControllers.GetBatches(ctx)
	})
}
