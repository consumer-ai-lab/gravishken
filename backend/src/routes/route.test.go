package route

import (
	"common/models/test"
	"fmt"
	"os"
	"server/src/controllers"
	"path/filepath"
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func TestRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	testRoute := route.Group("/test")
	// testRoute.Use(middleware.AdminJWTAuthMiddleware(allControllers.UserCollection))

	testRoute.POST("/add_test", func(ctx *gin.Context) {

		if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
			ctx.JSON(400, gin.H{"error": "File too large"})
			return
		}

		testType := ctx.Request.FormValue("type")
		duration := ctx.Request.FormValue("duration")
		typingText := ctx.Request.FormValue("typingText")

		fmt.Println("testType: ", testType)
		fmt.Println("duration: ", duration)
		fmt.Println("typing text: ", typingText)

		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}

		testModel := test.Test{
			Type:       test.TestType(testType),
			Duration:   durationInt,
			TypingText: typingText,
		}

		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			if err == http.ErrMissingFile {
				
				fmt.Println("No file uploaded, continuing without file")
			} else {
				ctx.JSON(400, gin.H{"error": "Error retrieving the file"})
				return
			}
		}
		if file != nil {
			sess, err := session.NewSession(&aws.Config{
				Region: aws.String("ap-south-1"),
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_S3_ACCESS_KEY"),
					os.Getenv("AWS_S3_ACCESS_KEY_SECRET"),
					"",
				),
			})
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to create AWS session"})
				return
			}

			uploader := s3manager.NewUploader(sess)

			filename := filepath.Base(header.Filename)
			result, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("collegeprojectbucket"),
				Key:    aws.String(filename),
				Body:   file,
			})

			if err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to upload file to S3"})
				return
			}

			testModel.File = result.Location
		}


		allControllers.AddTestToDB(ctx, &testModel);

		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to add test to database"})
			return
		}

		ctx.JSON(200, gin.H{"message": "Test added successfully", "test": testModel})
	})

	testRoute.GET("/get_question_paper/:batch_name", func(ctx *gin.Context) {

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

	testRoute.GET("/test_types", func(ctx *gin.Context) {
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

	testRoute.GET("/get_all_tests", func(ctx *gin.Context) {
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
