package controllers

import (
	"server/src/helper"
	User "server/src/models/user"

	"github.com/gin-gonic/gin"
)

func (this *Class) UserLoginHandler(ctx *gin.Context, userModel *User.UserLoginRequest) {
	userCollection := this.UserCollection
	response, err := helper.UserLogin(userCollection, userModel)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in User Login",
			"error": err,
		})
		return
	}


	ctx.JSON(200, gin.H{
		"message": "Admin Login route here",
		"response": response,
	})
}





