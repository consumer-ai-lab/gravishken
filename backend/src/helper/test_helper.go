package helper

import (
	"server/src/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestionPaper(Collection *mongo.Collection, password string) (types.ModelInterface, error) {
	var questionPaper types.ModelInterface

	err := Collection.FindOne(context.TODO(), bson.M{"password": password}).Decode(&questionPaper)

	if err != nil {
		return nil, err
	}
	return questionPaper, nil
}


func GetQuestionPaperByBatchNumber(Collection *mongo.Collection, batchNumber string) (types.ModelInterface, error) {
	var questionPaper types.ModelInterface

	err := Collection.FindOne(context.TODO(), bson.M{"batch": batchNumber}).Decode(&questionPaper)

	if err != nil {
		return nil, err
	}
	return questionPaper, nil
}