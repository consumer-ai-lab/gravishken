package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoutes(allControllers *controllers.Class, route *gin.Engine){
	testRoute := route.Group("/test")

	testRoute.GET("/get_question_paper/:password", func(ctx *gin.Context) {
		
		password := ctx.Param("password")

		questionPaper, err := allControllers.GetQuestionPaperHandler(ctx, password)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error while fetching question paper",
				"error": err,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Question paper fetched successfully",
			"questionPaper": questionPaper,
		})
	})
}