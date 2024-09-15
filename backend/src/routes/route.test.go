package route

import (
	"fmt"
	"strconv"
	"server/src/controllers"
	Test "common/models/test"
	"github.com/gin-gonic/gin"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
	// "path/filepath"
	// "context"
	// "go.mongodb.org/mongo-driver/bson"
)

func TestRoutes(allControllers *controllers.ControllerClass, route *gin.Engine){
	testRoute := route.Group("/test")

	testRoute.POST("/add_test",func(ctx *gin.Context){

		if err:=ctx.Request.ParseMultipartForm(10<<20);err!=nil {
			ctx.JSON(400,gin.H{"error":"File too large"})
			return
		}

		testType := ctx.Request.FormValue("type")
		duration := ctx.Request.FormValue("duration")
		typingText := ctx.Request.FormValue("typingText")

		fmt.Println("textType: ",testType)
		fmt.Println("duration: ",duration)
		fmt.Println("typing text: ",typingText)

		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}

		testModel := Test.Test{
			Type:       Test.TestType(testType),
			Duration:   durationInt,
			TypingText: typingText,
		}

		// file, header, err := ctx.Request.FormFile("file")

		if err == nil {
			//TODO: Handle Image to AWS
			// sess, err := session.NewSession(&aws.Config{
			// 	Region: aws.String("us-east-2"),
			// })
			// if err != nil {
			// 	ctx.JSON(500, gin.H{"error": "Failed to create AWS session"})
			// 	return
			// }

			// uploader := s3manager.NewUploader(sess)

			// filename := filepath.Base(header.Filename)
			// result, err := uploader.Upload(&s3manager.UploadInput{
			// 	Bucket: aws.String("wclsubmission"),
			// 	Key:    aws.String(filename),
			// 	Body:   file,
			// })
			
			// if err != nil {
			// 	ctx.JSON(500, gin.H{"error": "Failed to upload file to S3"})
			// 	return
			// }

			// testModel.File = result.Location
		}

		allControllers.AddTestToDB(ctx,&testModel)

		ctx.JSON(200, gin.H{"message": "Test added successfully", "test": testModel})
	})

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