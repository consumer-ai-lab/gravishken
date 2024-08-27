package helper

import (
	"context"
	"fmt"
	TEST "common/models/test"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestionPaper(Collection *mongo.Collection, password string) (TEST.Test, error) {
	var questionPaper TEST.Test

	err := Collection.FindOne(context.TODO(), bson.M{"password": password}).Decode(&questionPaper)

	if err != nil {
		return TEST.Test{}, err
	}
	return questionPaper, nil
}


func GetQuestionPaperByBatchNumber(Collection *mongo.Collection, batchNumber string) (TEST.Test, error) {
	var questionPaper TEST.Test

	fmt.Println("Inside getQuestionPaperByBatchNumber and batch number: ", batchNumber)

	err := Collection.FindOne(context.TODO(), bson.M{"batch": batchNumber}).Decode(&questionPaper)

	if err != nil {
		return TEST.Test{}, err
	}
	return questionPaper, nil
}