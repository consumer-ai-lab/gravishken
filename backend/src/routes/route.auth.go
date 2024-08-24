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

	allControllers := controllers.ControllerClass{
		Client:          db,
		AdminCollection: adminCollection,
		UserCollection:  userCollection,
		TestCollection:  testCollection,
		BatchCollection: batchCollection,
	}

	AdminRoutes(&allControllers, route)
	UserRoutes(&allControllers, route)
	BatchRoutes(&allControllers, route)
	TestRoutes(&allControllers, route)
}

