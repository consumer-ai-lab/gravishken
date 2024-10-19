package main

import (
	"common"
	"context"
	"log"

	// "encoding/csv"
	"fmt"
	// "io"
	// "log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func AdminRoutes(allControllers *Database, route *gin.Engine) {
	unauthenticatedAdminRoutes := route.Group("/admin")

	unauthenticatedAdminRoutes.POST("/register", func(ctx *gin.Context) {
		var adminModel common.Admin
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminRegisterHandler(ctx, &adminModel)
	})

	unauthenticatedAdminRoutes.POST("/login", func(ctx *gin.Context) {
		var adminModel common.Admin
		fmt.Println("Admin login route")
		if err := ctx.ShouldBindJSON(&adminModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.AdminLoginHandler(ctx, &adminModel)
	})

	authenticatedAdminRoutes := route.Group("/admin")
	authenticatedAdminRoutes.Use(AdminJWTAuthMiddleware(allControllers.AdminCollection))

	// If not authenticated, it will give 401 from the middleware
	authenticatedAdminRoutes.GET("/auth-status", func(ctx *gin.Context) {
		anyclaims, ok := ctx.Get("claims")
		if !ok {
			ctx.JSON(200, gin.H{
				"isAuthenticated": false,
				"error":           "Error fetching admin info",
			})
			return
		}
		claims := anyclaims.(*Claims)

		var adminInfo common.Admin
		err := allControllers.AdminCollection.FindOne(context.TODO(), bson.M{"username": claims.Username}).Decode(&adminInfo)
		if err != nil {
			ctx.JSON(200, gin.H{
				"isAuthenticated": false,
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

	// authenticatedAdminRoutes.POST("/add_users_from_csv", func(ctx *gin.Context) {
	// 	file, _, err := ctx.Request.FormFile("file")
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
	// 		return
	// 	}
	// 	defer file.Close()

	// 	reader := csv.NewReader(file)
	// 	var users []common.User

	// 	// Skip the header row
	// 	if _, err := reader.Read(); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format"})
	// 		return
	// 	}

	// 	for {
	// 		record, err := reader.Read()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error reading CSV"})
	// 			return
	// 		}

	// 		if len(record) != 5 {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format"})
	// 			return
	// 		}

	// 		user := common.User{
	// 			ID:           primitive.NewObjectID(),
	// 			Username:     record[0],
	// 			Password:     record[1],
	// 			TestPassword: record[2],
	// 			Batch:        record[3],
	// 			Tests:        common.UserSubmission{},
	// 		}
	// 		users = append(users, user)
	// 	}

	// 	// Insert users into the database
	// 	userInterfaces := make([]interface{}, len(users))
	// 	for i, u := range users {
	// 		userInterfaces[i] = u
	// 	}

	// 	// Insert users into the database
	// 	insertedResult, err := allControllers.UserCollection.InsertMany(context.Background(), userInterfaces)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert users"})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": fmt.Sprintf("Successfully added %d users", len(insertedResult.InsertedIDs)),
	// 	})
	// })

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
		testObjectIDs := make([]string, 0, len(batchData.SelectedTests))
		for _, testID := range batchData.SelectedTests {
			objectID, err := primitive.ObjectIDFromHex(testID)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid test ID format"})
				return
			}
			testObjectIDs = append(testObjectIDs, objectID.String())
		}

		newBatch := common.Batch{
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

		testModel := common.Test{
			Type:       common.TestType(testType),
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

		allControllers.AddTestToDB(ctx, &testModel)

		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to add test to database"})
			return
		}

		ctx.JSON(200, gin.H{"message": "Test added successfully", "test": testModel})
	})

	// authenticatedAdminRoutes.POST("/update_user_data", func(ctx *gin.Context) {
	// 	var userUpdateRequest common.UserUpdateRequest

	// 	if err := ctx.ShouldBindJSON(&userUpdateRequest); err != nil {
	// 		ctx.JSON(500, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}

	// 	allControllers.UpdateUserData(ctx, &userUpdateRequest)

	// })

	// authenticatedAdminRoutes.POST("/increase_test_time", func(ctx *gin.Context) {
	// 	var requestData struct {
	// 		Param          string   `json:"param"`
	// 		Username       []string `json:"username"`
	// 		TimeToIncrease int64    `json:"time_to_increase"`
	// 	}

	// 	// Bind JSON body to requestData struct
	// 	if err := ctx.ShouldBindJSON(&requestData); err != nil {
	// 		ctx.JSON(400, gin.H{
	// 			"message": "Invalid request body",
	// 			"error":   err.Error(),
	// 		})
	// 		return
	// 	}

	// 	allControllers.Increase_Time(ctx, requestData.Param, requestData.Username, requestData.TimeToIncrease)

	// })

	// authenticatedAdminRoutes.POST("/get_batchwise_data", func(ctx *gin.Context) {
	// 	var batchData struct {
	// 		Param       string `json:"param"`
	// 		BatchNumber string `json:"batchNumber"`
	// 		Ranges      []int  `json:"ranges"`
	// 	}

	// 	if err := ctx.ShouldBindJSON(&batchData); err != nil {
	// 		ctx.JSON(500, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}

	// 	allControllers.GetBatchWiseData(ctx, batchData.Param, batchData.BatchNumber, batchData.Ranges)

	// })

	// authenticatedAdminRoutes.POST("/set_user_data", func(ctx *gin.Context) {
	// 	var userRequest struct {
	// 		Username         string `json:"username"`
	// 		Param            string `json:"param"`
	// 		From             int    `json:"from"`
	// 		To               int    `json:"to"`
	// 		ResultDownloaded bool   `json:"resultDownloaded"`
	// 	}

	// 	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
	// 		ctx.JSON(500, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}

	// 	allControllers.SetUserData(ctx, userRequest.Param, &common.UserBatchRequestData{
	// 		From:             userRequest.From,
	// 		To:               userRequest.To,
	// 		ResultDownloaded: userRequest.ResultDownloaded,
	// 	}, userRequest.Username)
	// })

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

func InitAuthRoutes(db *Database, route *gin.Engine) {
	AdminRoutes(db, route)
	UserRoutes(db, route)
	BatchRoutes(db, route)
	TestRoutes(db, route)
}

func BatchRoutes(allControllers *Database, route *gin.Engine) {
	batchRoute := route.Group("/batch")
	// batchRoute.Use(UserJWTAuthMiddleware(allControllers.UserCollection))

	batchRoute.GET("/get_batches", func(ctx *gin.Context) {
		allControllers.GetBatches(ctx)
	})
}

func TestRoutes(allControllers *Database, route *gin.Engine) {
	unauthenticatedTestRoute := route.Group("/test")
	authenticatedTestRoute := route.Group("/test")
	authenticatedTestRoute.Use(UserJWTAuthMiddleware(allControllers.UserCollection))

	authenticatedTestRoute.GET("/get_question_paper/:batch_name", func(ctx *gin.Context) {

		batch_name := ctx.Param("batch_name")
		log.Println(batch_name)

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
			string(common.TypingTest),
			string(common.DocxTest),
			string(common.ExcelTest),
			string(common.WordTest),
		}

		ctx.JSON(200, gin.H{
			"message":   "Test types fetched successfully",
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

func UserRoutes(allControllers *Database, route *gin.Engine) {
	userRoute := route.Group("/user")

	userRoute.POST("/login", func(ctx *gin.Context) {
		var userModel common.TUserLoginRequest
		if err := ctx.ShouldBindJSON(&userModel); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		allControllers.UserLoginHandler(ctx, &userModel)
	})

	authenticated := userRoute.Group("/")

	authenticated.Use(AdminJWTAuthMiddleware(allControllers.UserCollection))

	authenticated.GET("/get_all_users", func(ctx *gin.Context) {
		var users []common.User
		cursor, err := allControllers.UserCollection.Find(context.Background(), bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &users); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"users": users})
	})

	authenticated.GET("/paginated_users", func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
		skip := (page - 1) * limit

		var users []common.User

		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
		cursor, err := allControllers.UserCollection.Find(context.Background(), bson.M{}, opts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &users); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		totalUsers, err := allControllers.UserCollection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total users"})
			return
		}

		totalPages := int(math.Ceil(float64(totalUsers) / float64(limit)))

		ctx.JSON(http.StatusOK, gin.H{
			"users":       users,
			"totalPages":  totalPages,
			"currentPage": page,
		})
	})

	// userRoute.PUT("/update_user", func(ctx *gin.Context) {
	// 	var updateRequest common.UserModelUpdateRequest
	// 	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}

	// 	log.Default().Println("User Request: ", updateRequest)

	// 	err := allControllers.UpdateUser(ctx, &updateRequest)

	// 	if err != nil {
	// 		ctx.JSON(500, gin.H{
	// 			"message": "Error in updating user!",
	// 			"error":   err,
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(200, gin.H{
	// 		"message": "user updated successfully!!",
	// 	})
	// })

	userRoute.DELETE("/delete_user", func(ctx *gin.Context) {
		var deleteRequest struct {
			UserId string `json:"userId"`
		}
		if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if deleteRequest.UserId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		err := allControllers.DeleteUser(ctx, deleteRequest.UserId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

}
