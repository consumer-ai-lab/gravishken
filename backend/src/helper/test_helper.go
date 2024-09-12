package helper

import (
	TEST "common/models/test"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestionPaper(Collection *mongo.Collection, password string) (TEST.BatchTests, error) {
	var questionPaper TEST.BatchTests

	err := Collection.FindOne(context.TODO(), bson.M{"password": password}).Decode(&questionPaper)

	if err != nil {
		return TEST.BatchTests{}, err
	}
	return questionPaper, nil
}

func GetQuestionPaperByBatchNumber(Collection *mongo.Collection, batchNumber string) (TEST.BatchTests, error) {
	var questionPaper TEST.BatchTests

	fmt.Println("Inside getQuestionPaperByBatchNumber and batch number: ", batchNumber)

	err := Collection.FindOne(context.TODO(), bson.M{"batch": batchNumber}).Decode(&questionPaper)

	if err != nil {
		return TEST.BatchTests{}, err
	}
	return questionPaper, nil
}
