package controllers

import (
	"common/models/admin"
	Test "common/models/test"
	User "common/models/user"
	"log"
	"server/src/helper"
	"server/src/utils"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) AdminLoginHandler(ctx *gin.Context, adminModel *admin.Admin) {
	adminCollection := this.AdminCollection
	response, err := helper.AdminLogin(adminCollection, adminModel)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in Admin Login",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message":  "Admin Login route here",
		"response": response,
	})
}

func (this *ControllerClass) AdminRegisterHandler(ctx *gin.Context, adminModel *admin.Admin) {
	adminCollection := this.AdminCollection
	err := helper.RegisterAdmin(adminCollection, adminModel)

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

func (this *ControllerClass) AdminChangePasswordHandler(ctx *gin.Context, adminModel *admin.AdminChangePassword) {
	adminCollection := this.AdminCollection
	err := helper.ChangePassword(adminCollection, adminModel)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in Admin Change Password",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Admin Change Password route here",
	})
}

func (this *ControllerClass) AddTestToDB(ctx *gin.Context, test *Test.Test) {
	testCollection := this.TestCollection
	err := helper.Add_Model_To_DB(testCollection, test)

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

func (this *ControllerClass) AddAllUsersBacthesToDb(ctx *gin.Context, filePath string) {
	userCollection := this.UserCollection
	testCollection := this.TestCollection

	csvData, unique_batches := utils.Read_CSV(filePath)

	// creating a map to store test passwords for each batch
	batch_passwords := make(map[string]string)

	log.Default().Println("Adding all batches to db")

	// Looping over all batches and finding test password for each batch and storing it in a map
	for batch, _ := range unique_batches {
		test_data, err := helper.GetQuestionPaperByBatchNumber(testCollection, batch)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error while fetching question paper",
				"error":   err,
			})
			return
		}

		batch_passwords[batch] = test_data.Password
	}

	// Looping over all user data fetched from reading csv file and adding them to db
	for _, data := range csvData {
		user := User.User{
			Name:         data["name"],
			Username:     data["roll_no"],
			Password:     data["password"],
			TestPassword: batch_passwords[data["slot"]],
			Batch:        data["slot"],
			Tests:        User.UserSubmission{},
		}

		helper.Add_Model_To_DB(userCollection, &user)
	}

	ctx.JSON(200, gin.H{
		"message": "All users added to db",
	})
}

func (this *ControllerClass) UpdateTypingTestText(ctx *gin.Context, typingTestText string, testPassword string) {
	testCollection := this.TestCollection

	err := helper.UpdateTypingTestText(testCollection, typingTestText, testPassword)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error while updating typing test text",
			"error":   err,
		})
	}
}
