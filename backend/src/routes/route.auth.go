package route

import (
	"fmt"
	"server/src/controllers"
	"server/src/models/admin"

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
}

func AdminRoutes(allControllers *controllers.Class, route *gin.Engine){
	adminRoute := route.Group("/admin")


	adminRoute.POST("/login", func(ctx *gin.Context) {
		var adminModel admin.Admin
		fmt.Println("Admin login route")
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminLoginHandler(ctx, &adminModel)
	})


	adminRoute.POST("/register", func(ctx *gin.Context) {
		var adminModel admin.Admin
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		
		allControllers.AdminRegisterHandler(ctx, &adminModel)
	})

	adminRoute.POST("/changepassword", func(ctx *gin.Context) {
		var adminModel admin.AdminChangePassword
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		
		allControllers.AdminChangePasswordHandler(ctx, &adminModel)
	})

	adminRoute.POST("/add_all_users", func(ctx *gin.Context) {
		var filePath string
		if err := ctx.ShouldBindJSON(&filePath); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		allControllers.AddAllUsers(ctx, filePath)
	})
}

func SampleHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Auth routes here",
	})
}
