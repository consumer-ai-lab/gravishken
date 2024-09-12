package controllers

import (
	Test "common/models/test"
	"server/src/helper"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

func (this *ControllerClass) GetQuestionPaperHandler(ctx *gin.Context, password string) (types.ModelInterface, error) {
	testCollection := this.TestCollection

	testModel, err := helper.GetQuestionPaper(testCollection, password)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error while fetching question paper"})
		return &Test.BatchTests{}, err
	}

	return &testModel, nil
}
