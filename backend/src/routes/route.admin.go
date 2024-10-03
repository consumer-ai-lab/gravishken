package route

import (
	"common/models/admin"
	"common/models/batch"
	"common/models/test"
	"common/models/user"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"net/http"
	"server/src/controllers"
	"server/src/middleware"
	"server/src/types"
	"os"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func AdminRoutes(allControllers *controllers.ControllerClass, route *gin.Engine) {
	unauthenticatedAdminRoutes := route.Group("/admin")

	unauthenticatedAdminRoutes.POST("/register", func(ctx *gin.Context) {
		var adminModel admin.Admin
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminRegisterHandler(ctx, &adminModel)
	})

	unauthenticatedAdminRoutes.POST("/login", func(ctx *gin.Context) {
		var adminModel admin.Admin
		fmt.Println("Admin login route")
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminLoginHandler(ctx, &adminModel)
	})

	authenticatedAdminRoutes := route.Group("/admin")
	authenticatedAdminRoutes.Use(middleware.AdminJWTAuthMiddleware(allControllers.AdminCollection))

	authenticatedAdminRoutes.GET("/auth-status", func(ctx *gin.Context) {
		anyclaims, ok := ctx.Get("claims")
		if !ok {
			ctx.JSON(500, gin.H{
				"isAuthenticated": true,
				"error":           "Error fetching admin info",
			})
			return
		}
		claims := anyclaims.(*types.Claims)

		var adminInfo admin.Admin
		err := allControllers.AdminCollection.FindOne(context.TODO(), bson.M{"username": claims.Username}).Decode(&adminInfo)
		if err != nil {
			ctx.JSON(500, gin.H{
				"isAuthenticated": true,
				"error":           "Error fetching admin info",
			})
			return
		}

		// Remove sensitive information
		adminInfo.Password = ""

		ctx.JSON(200, gin.H{
			"isAuthenticated": true,
			"adminInfo":       adminInfo,
		})
	})

	authenticatedAdminRoutes.POST("/add_all_users", func(ctx *gin.Context) {
		var FilePathRequest struct {
			FilePath string `json:"filePath" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&FilePathRequest); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		allControllers.AddAllUsersBacthesToDb(ctx, FilePathRequest.FilePath)
	})

	authenticatedAdminRoutes.POST("/add_batch", func(ctx *gin.Context) {
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

	authenticatedAdminRoutes.POST("/add_test", func(ctx *gin.Context) {

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

	authenticatedAdminRoutes.POST("/update_user_data", func(ctx *gin.Context) {
		var userUpdateRequest user.UserUpdateRequest

		if err := ctx.ShouldBindJSON(&userUpdateRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UpdateUserData(ctx, &userUpdateRequest)

	})

	authenticatedAdminRoutes.POST("/increase_test_time", func(ctx *gin.Context) {
		var requestData struct {
			Param          string   `json:"param"`
			Username       []string `json:"username"`
			TimeToIncrease int64    `json:"time_to_increase"`
		}

		// Bind JSON body to requestData struct
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			return
		}

		allControllers.Increase_Time(ctx, requestData.Param, requestData.Username, requestData.TimeToIncrease)

	})

	authenticatedAdminRoutes.POST("/get_batchwise_data", func(ctx *gin.Context) {
		var batchData struct {
			Param       string `json:"param"`
			BatchNumber string `json:"batchNumber"`
			Ranges      []int  `json:"ranges"`
		}

		if err := ctx.ShouldBindJSON(&batchData); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.GetBatchWiseData(ctx, batchData.Param, batchData.BatchNumber, batchData.Ranges)

	})

	authenticatedAdminRoutes.POST("/set_user_data", func(ctx *gin.Context) {
		var userRequest struct {
			Username         string `json:"username"`
			Param            string `json:"param"`
			From             int    `json:"from"`
			To               int    `json:"to"`
			ResultDownloaded bool   `json:"resultDownloaded"`
		}

		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.SetUserData(ctx, userRequest.Param, &user.UserBatchRequestData{
			From:             userRequest.From,
			To:               userRequest.To,
			ResultDownloaded: userRequest.ResultDownloaded,
		}, userRequest.Username)
	})

	authenticatedAdminRoutes.POST("/update_typing_test_text", func(ctx *gin.Context) {
		var UpdateTypingTestTextRequest struct {
			TypingTestText string `json:"typingTestText"`
			TestPassword   string `json:"testPassword"`
		}

		if err := ctx.ShouldBindJSON(&UpdateTypingTestTextRequest); err != nil {
			ctx.JSON(500, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UpdateTypingTestText(ctx, UpdateTypingTestTextRequest.TypingTestText, UpdateTypingTestTextRequest.TestPassword)
	})
}
