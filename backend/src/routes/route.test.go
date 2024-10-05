package route

import (
	"common/models/test"
	"server/src/controllers"
	"server/src/middleware"
	"github.com/gin-gonic/gin"

)

func TestRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	unauthenticatedTestRoute := route.Group("/test")
	authenticatedTestRoute := route.Group("/test")
	authenticatedTestRoute.Use(middleware.AdminJWTAuthMiddleware(allControllers.UserCollection))

	authenticatedTestRoute.GET("/get_question_paper/:batch_name", func(ctx *gin.Context) {

		batch_name := ctx.Param("batch_name")

		questionPaper, err := allControllers.GetQuestionPaperHandler(ctx, batch_name)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error while fetching question paper",
				"error":   err,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message":       "Question paper fetched successfully",
			"questionPaper": questionPaper,
		})
	})

	unauthenticatedTestRoute.GET("/test_types", func(ctx *gin.Context) {
        testTypes := []string{
            string(test.TypingTest),
            string(test.DocxTest),
            string(test.ExcelTest),
            string(test.WordTest),
        }

        ctx.JSON(200, gin.H{
            "message": "Test types fetched successfully",
            "testTypes": testTypes,
        })
    })

	unauthenticatedTestRoute.GET("/get_all_tests", func(ctx *gin.Context) {
		tests, err := allControllers.GetAllTests(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error while fetching tests",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Tests fetched successfully",
			"tests":   tests,
		})
	})
}
