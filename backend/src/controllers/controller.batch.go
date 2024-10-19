package controllers

import (
	"common/models/batch"
	"server/src/helper"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) AddBatchToDB(ctx *gin.Context, batchData *batch.Batch) {
	testCollection := this.BatchCollection

	err := helper.Add_Model_To_DB(testCollection, batchData)

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

func (this *ControllerClass) GetBatches(ctx *gin.Context) {
	testCollection := this.BatchCollection

	batchData, err := helper.Get_All_Models(testCollection, &batch.Batch{})

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
