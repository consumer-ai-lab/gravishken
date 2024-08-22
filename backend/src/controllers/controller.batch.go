package controllers

import (
	"server/src/utils"

	"github.com/gin-gonic/gin"
)


func (this *Class) AddBatchToDB(ctx *gin.Context, filePath string) {
	testCollection := this.TestCollection
	err := utils.Add_CSVData_To_DB(testCollection, filePath)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error while adding batch to DB",
			"error": err,
		})
	}
	ctx.JSON(200, gin.H{
		"message": "Batch added to DB successfully",
	})
}