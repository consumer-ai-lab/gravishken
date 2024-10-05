package controllers

import (
	"common/models/test"
	// "common/models/batch"
	"server/src/helper"
	"server/src/types"

	"github.com/gin-gonic/gin"
	"context"


	"go.mongodb.org/mongo-driver/bson"
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


func (c *ControllerClass) GetAllTests(ctx *gin.Context) ([]test.Test, error) {
	var tests []test.Test

	cursor, err := c.TestCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &tests); err != nil {
		return nil, err
	}

	return tests, nil
}




