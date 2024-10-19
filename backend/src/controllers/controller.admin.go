package controllers

import (
	"common/models/admin"
	Test "common/models/test"
	// User "common/models/user"
	// "log"
	"server/src/helper"
	// "server/src/utils"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) AdminLoginHandler(ctx *gin.Context, adminModel *admin.Admin) {
	adminCollection := this.AdminCollection
	token, err := helper.AdminLogin(adminCollection, adminModel)

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

func (this *ControllerClass) AdminChangePassword(ctx *gin.Context) {
	ctx.JSON(501, gin.H{
		"message": "This route is not needed",
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

func (this *ControllerClass) UpdateTypingTestText(ctx *gin.Context, typingTestText string, testID string) {
	testCollection := this.TestCollection

	err := helper.UpdateTypingTestText(testCollection, testID, typingTestText)
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
