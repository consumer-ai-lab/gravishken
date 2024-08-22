package controllers

import (
	"server/src/helper"
	Test "server/src/models/test"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

func (this *Class) GetQuestionPaperHandler(ctx *gin.Context, password string)  (types.ModelInterface, error) {
	testCollection := this.TestCollection
	testModel, err := helper.GetQuestionPaper(testCollection, password)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error while fetching question paper"})
		return &Test.Test{}, err
	}

	return testModel, nil
}