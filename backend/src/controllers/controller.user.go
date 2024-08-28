package controllers

import (
	User "common/models/user"
	"server/src/helper"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) UserLoginHandler(ctx *gin.Context, userModel *User.UserLoginRequest) {
	userCollection := this.UserCollection
	response, err := helper.UserLogin(userCollection, userModel)

	if err != nil {
		ctx.JSON(401, gin.H{
			"message": "Error in User Login",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message":  "Admin Login route here",
		"response": response,
	})
}

func (this *ControllerClass) UpdateUserData(ctx *gin.Context, userUpdateRequest *User.UserUpdateRequest) {
	userCollection := this.UserCollection
	err := helper.UpdateUserData(userCollection, userUpdateRequest)

	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error in updating user data",
			"error":   err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "User data updated successfully",
	})
}

func (this *ControllerClass) Increase_Time(ctx *gin.Context, param string, username []string, time_to_increase int64) {
	userCollection := this.UserCollection

	if len(username) == 0 {
		ctx.JSON(500, gin.H{
			"message": "Empty username",
		})
		return
	}

	if len(username) > 1 {
		param = "batch"
	}

	switch param {
	case "user":
		err := helper.UpdateUserTestTime(userCollection, username[0], time_to_increase)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in increasing time",
				"error":   err,
			})
		}
		ctx.JSON(200, gin.H{
			"message": "Time increased successfully",
		})

	case "batch":

		err := helper.UpdateBatchTestTime(userCollection, username, time_to_increase)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in increasing time",
				"error":   err,
			})
		}
		ctx.JSON(200, gin.H{
			"message": "Time increased successfully",
		})

	default:
		ctx.JSON(500, gin.H{
			"message": "Invalid parameter",
		})
	}

}

func (this *ControllerClass) GetBatchWiseData(ctx *gin.Context, param string, BatchNumber string, Ranges []int) {
	userCollection := this.UserCollection

	switch param {
	case "batch":
		result, err := helper.GetBatchWiseList(userCollection, BatchNumber)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in fetching batch wise data",
				"error":   err,
			})
		}
		ctx.JSON(200, gin.H{
			"message": "Batch wise data fetched successfully",
			"data":    result,
		})

	case "roll":
		From := Ranges[0]
		To := Ranges[1]
		result, err := helper.GetBatchWiseListRoll(userCollection, BatchNumber, From, To)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in fetching batch wise data",
				"error":   err,
			})
		}

		ctx.JSON(200, gin.H{
			"message": "Batch wise data fetched successfully",
			"data":    result,
		})

	case "frontend":
		result, err := helper.GetBatchDataForFrontend(userCollection, BatchNumber)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in fetching batch wise data",
				"error":   err,
			})
		}

		ctx.JSON(200, gin.H{
			"message": "Batch wise data fetched successfully",
			"data":    result,
		})

	default:
		ctx.JSON(500, gin.H{
			"message": "Invalid parameter",
		})
	}
}

func (this *ControllerClass) SetUserData(ctx *gin.Context, param string, userRequest *User.UserBatchRequestData, Username string) {
	userCollection := this.UserCollection

	switch param {
	case "download":
		err := helper.SetUserResultToDownloaded(userCollection, userRequest)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in setting user data",
				"error":   err,
			})
		}

		ctx.JSON(200, gin.H{
			"message": "User data set successfully",
		})

	case "reset":
		err := helper.ResetUserData(userCollection, Username)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error in resetting user data",
				"error":   err,
			})
		}

		ctx.JSON(200, gin.H{
			"message": "User data reset successfully",
		})

	default:
		ctx.JSON(500, gin.H{
			"message": "Invalid parameter",
		})
	}

}
