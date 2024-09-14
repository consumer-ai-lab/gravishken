package helper

import (
	"common/models/batch"
	"common/models/test"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTestsByBatch(batchCollection *mongo.Collection, testCollection *mongo.Collection, batchName string) ([]test.Test, error) {
	var batchDoc batch.Batch
	err := batchCollection.FindOne(context.TODO(), bson.M{"name": batchName}).Decode(&batchDoc)
	if err != nil {
		return nil, fmt.Errorf("error finding batch: %v", err)
	}

	var tests []test.Test
	cursor, err := testCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": batchDoc.Tests}})
	if err != nil {
		return nil, fmt.Errorf("error finding tests: %v", err)
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &tests)
	if err != nil {
		return nil, fmt.Errorf("error decoding tests: %v", err)
	}

	return tests, nil
}

func GetTestByID(testCollection *mongo.Collection, testID primitive.ObjectID) (*test.Test, error) {
	var testDoc test.Test
	err := testCollection.FindOne(context.TODO(), bson.M{"_id": testID}).Decode(&testDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("error finding test: %v", err)
	}
	return &testDoc, nil
}

