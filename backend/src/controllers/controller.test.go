package controllers

import (
	// "common/models/test"
	// "common/models/batch"
	"server/src/helper"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) GetQuestionPaperHandler(ctx *gin.Context, batchName string) ([]types.ModelInterface, error) {
	batchCollection := this.BatchCollection
	testCollection := this.TestCollection

	tests, err := helper.GetTestsByBatch(batchCollection, testCollection, batchName)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error while fetching question papers"})
		return nil, err
	}

	// Convert []test.Test to []types.ModelInterface
	var modelTests []types.ModelInterface
	for _, t := range tests {
		modelTests = append(modelTests, &t)
	}

	return modelTests, nil
}
