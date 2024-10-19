package route

import (
	"server/src/controllers"
	// "server/src/middleware"

	"github.com/gin-gonic/gin"
)

func BatchRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	batchRoute := route.Group("/batch")
	// batchRoute.Use(middleware.UserJWTAuthMiddleware(allControllers.UserCollection))

	batchRoute.GET("/get_batches", func(ctx *gin.Context) {
		allControllers.GetBatches(ctx)
	})
}
