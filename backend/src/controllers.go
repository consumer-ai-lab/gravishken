package main

import (
	"common"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client          *mongo.Client
	AdminCollection *mongo.Collection
	UserCollection  *mongo.Collection
	TestCollection  *mongo.Collection
	BatchCollection *mongo.Collection
}

func connectDatabase() (*Database, error) {
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")

	adminCollection := client.Database("GRAVTEST").Collection("Admin")
	userCollection := client.Database("GRAVTEST").Collection("Users")
	testCollection := client.Database("GRAVTEST").Collection("Tests")
	batchCollection := client.Database("GRAVTEST").Collection("Batch")

	db := Database{
		Client:          client,
		AdminCollection: adminCollection,
		UserCollection:  userCollection,
		TestCollection:  testCollection,
		BatchCollection: batchCollection,
	}
	return &db, nil
}

func (this *Database) GetQuestionPaperHandler(ctx *gin.Context, batchName string) ([]ModelInterface, error) {
	batchCollection := this.BatchCollection
	testCollection := this.TestCollection

	tests, err := GetTestsByBatch(batchCollection, testCollection, batchName)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error while fetching question papers"})
		return nil, err
	}

	// Convert []test.Test to []types.ModelInterface
	var modelTests []ModelInterface
	for _, t := range tests {
		modelTests = append(modelTests, &t)
	}

	return modelTests, nil
}

func (c *Database) GetAllTests(ctx *gin.Context) ([]common.Test, error) {
	var tests []common.Test

	cursor, err := c.TestCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &tests); err != nil {
		return nil, err
	}

	return tests, nil
}

func (this *Database) AdminLoginHandler(ctx *gin.Context, adminModel *common.Admin) {
	adminCollection := this.AdminCollection
	token, err := AdminLogin(adminCollection, adminModel)

	if err != nil {
		ctx.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set the token in a cookie
	ctx.SetCookie("auth_token", token, 3600*48, "/", "", false, true)

	ctx.JSON(200, gin.H{
		"message": "Admin logged in successfully",
	})
}

func (this *Database) AdminRegisterHandler(ctx *gin.Context, adminModel *common.Admin) {
	adminCollection := this.AdminCollection
	err := RegisterAdmin(adminCollection, adminModel)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in Admin Register",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Admin Register route here",
	})
}

func (this *Database) AdminChangePassword(ctx *gin.Context) {
	ctx.JSON(501, gin.H{
		"message": "This route is not needed",
	})
}

func (this *Database) AddTestToDB(ctx *gin.Context, test *common.Test) {
	testCollection := this.TestCollection
	err := Add_Model_To_DB(testCollection, test)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error while adding test to db",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Test added to db",
	})
}

func (this *Database) UpdateTypingTestText(ctx *gin.Context, typingTestText string, testID string) {
	testCollection := this.TestCollection

	err := UpdateTypingTestText(testCollection, testID, typingTestText)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error while updating typing test text",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Typing test text updated successfully",
	})
}

func (this *Database) AddBatchToDB(ctx *gin.Context, batchData *common.Batch) {
	testCollection := this.BatchCollection

	err := Add_Model_To_DB(testCollection, batchData)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in adding batch data",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Batch data added successfully",
	})
}

func (this *Database) GetBatches(ctx *gin.Context) {
	testCollection := this.BatchCollection

	batchData, err := Get_All_Models(testCollection, &common.Batch{})

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in fetching batch data",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Batch data fetched successfully",
		"data":    batchData,
	})
}

func (this *Database) UserLoginHandler(ctx *gin.Context, userModel *common.TUserLoginRequest) {
	userCollection := this.UserCollection
	response, err := UserLogin(userCollection, userModel)
	user, _ := common.FindByUsername(userCollection, userModel.Username)

	if err != nil {
		ctx.JSON(401, gin.H{
			"message": "Error in User Login",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, common.UserLoginResponse{
		Jwt:  response,
		User: *user,
	})
}

// func (this *ControllerClass) UpdateUserData(ctx *gin.Context, userUpdateRequest *common.UserUpdateRequest) {
// 	userCollection := this.UserCollection
// 	err := UpdateUserData(userCollection, userUpdateRequest)

// 	if err != nil {
// 		ctx.JSON(500, gin.H{
// 			"message": "Error in updating user data",
// 			"error":   err,
// 		})
// 		return
// 	}

// 	ctx.JSON(200, gin.H{
// 		"message": "User data updated successfully",
// 	})
// }

// func (this *ControllerClass) Increase_Time(ctx *gin.Context, param string, username []string, time_to_increase int64) {
// 	userCollection := this.UserCollection

// 	if len(username) == 0 {
// 		ctx.JSON(500, gin.H{
// 			"message": "Empty username",
// 		})
// 		return
// 	}

// 	if len(username) > 1 {
// 		param = "batch"
// 	}

// 	switch param {
// 	case "user":
// 		err := UpdateUserTestTime(userCollection, username[0], time_to_increase)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in increasing time",
// 				"error":   err,
// 			})
// 		}
// 		ctx.JSON(200, gin.H{
// 			"message": "Time increased successfully",
// 		})

// 	case "batch":

// 		err := UpdateBatchTestTime(userCollection, username, time_to_increase)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in increasing time",
// 				"error":   err,
// 			})
// 		}
// 		ctx.JSON(200, gin.H{
// 			"message": "Time increased successfully",
// 		})

// 	default:
// 		ctx.JSON(500, gin.H{
// 			"message": "Invalid parameter",
// 		})
// 	}

// }

// func (this *ControllerClass) GetBatchWiseData(ctx *gin.Context, param string, BatchNumber string, Ranges []int) {
// 	userCollection := this.UserCollection

// 	switch param {
// 	case "batch":
// 		result, err := GetBatchWiseList(userCollection, BatchNumber)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in fetching batch wise data",
// 				"error":   err,
// 			})
// 		}
// 		ctx.JSON(200, gin.H{
// 			"message": "Batch wise data fetched successfully",
// 			"data":    result,
// 		})

// 	case "roll":
// 		From := Ranges[0]
// 		To := Ranges[1]
// 		result, err := GetBatchWiseListRoll(userCollection, BatchNumber, From, To)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in fetching batch wise data",
// 				"error":   err,
// 			})
// 		}

// 		ctx.JSON(200, gin.H{
// 			"message": "Batch wise data fetched successfully",
// 			"data":    result,
// 		})

// 	case "frontend":
// 		result, err := GetBatchDataForFrontend(userCollection, BatchNumber)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in fetching batch wise data",
// 				"error":   err,
// 			})
// 		}

// 		ctx.JSON(200, gin.H{
// 			"message": "Batch wise data fetched successfully",
// 			"data":    result,
// 		})

// 	default:
// 		ctx.JSON(500, gin.H{
// 			"message": "Invalid parameter",
// 		})
// 	}
// }

// func (this *ControllerClass) SetUserData(ctx *gin.Context, param string, userRequest *common.UserBatchRequestData, Username string) {
// 	userCollection := this.UserCollection

// 	switch param {
// 	case "download":
// 		err := SetUserResultToDownloaded(userCollection, userRequest)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in setting user data",
// 				"error":   err,
// 			})
// 		}

// 		ctx.JSON(200, gin.H{
// 			"message": "User data set successfully",
// 		})

// 	case "reset":
// 		err := ResetUserData(userCollection, Username)
// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error in resetting user data",
// 				"error":   err,
// 			})
// 		}

// 		ctx.JSON(200, gin.H{
// 			"message": "User data reset successfully",
// 		})

// 	default:
// 		ctx.JSON(500, gin.H{
// 			"message": "Invalid parameter",
// 		})
// 	}

// }

// func (self *ControllerClass) UpdateUser(ctx *gin.Context, userRequest *common.UserModelUpdateRequest) error {
// 	userCollection := self.UserCollection

// 	err := UpdateUser(userCollection, userRequest)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (self *Database) DeleteUser(ctx *gin.Context, userId string) error {
	userCollection := self.UserCollection

	err := Delete_Model_By_ID(userCollection, userId)

	if err != nil {
		return err
	}

	return nil
}
