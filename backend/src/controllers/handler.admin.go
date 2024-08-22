package controllers

import (
	"server/src/helper"
	"server/src/models/admin"

	"github.com/gin-gonic/gin"
)

func (this *Class) AdminLoginHandler(ctx *gin.Context, adminModel *admin.Admin) {
	adminCollection := this.AdminCollection
	response, err := helper.AdminLogin(adminCollection, adminModel)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in Admin Login",
			"error": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Admin Login route here",
		"response": response,
	})
}

func (this *Class) AdminRegisterHandler(ctx *gin.Context, adminModel *admin.Admin) {
	adminCollection := this.AdminCollection
	response := helper.RegisterAdmin(adminCollection, adminModel)
	ctx.JSON(200, gin.H{
		"message": "Admin Register route here",
		"response": response,
	})
}

func (this *Class) AdminChangePasswordHandler(ctx *gin.Context, adminModel *admin.AdminChangePassword) {
	adminCollection := this.AdminCollection
	response := helper.ChangePassword(adminCollection, adminModel)
	ctx.JSON(200, gin.H{
		"message": "Admin Change Password route here",
		"response": response,
	})
}

func (this *Class) AddAllUsers(ctx *gin.Context, filePath string){

}
