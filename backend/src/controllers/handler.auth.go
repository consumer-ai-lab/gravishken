package controllers

import "github.com/gin-gonic/gin"

func (this *Class) AdminLoginHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Admin Login route here",
	})
}
