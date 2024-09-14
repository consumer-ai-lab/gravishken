package route

import (
	"server/src/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoutes(allControllers *controllers.ControllerClass, route *gin.Engine){
	testRoute := route.Group("/test")

	testRoute.GET("/get_question_paper/:batch_name", func(ctx *gin.Context) {
		
		batch_name := ctx.Param("batch_name")

		questionPaper, err := allControllers.GetQuestionPaperHandler(ctx, batch_name)
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