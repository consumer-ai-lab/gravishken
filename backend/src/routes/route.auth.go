package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitAuthRoutes(db *mongo.Client, route *gin.Engine) {
	adminCollection := db.Database("GRAVTEST").Collection("Admin")
	userCollection := db.Database("GRAVTEST").Collection("Users")
	testCollection := db.Database("GRAVTEST").Collection("Tests")
	batchCollection := db.Database("GRAVTEST").Collection("Batch")

	allControllers := controllers.Class{
		Client: db,
		AdminCollection: adminCollection,
		UserCollection: userCollection,
		TestCollection: testCollection,
		BatchCollection: batchCollection,
	}

	AdminRoutes(&allControllers, route)
	BatchRoutes(&allControllers, route)
	TestRoutes(&allControllers, route)
}


func SampleHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Auth routes here",
	})
}
