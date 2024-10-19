package helper

import (
	"context"
	"server/src/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBatchByBatchNumber(Collection *mongo.Collection, batchNumber string) (types.ModelInterface, error) {
	var batch types.ModelInterface

	err := Collection.FindOne(context.TODO(), bson.M{"batch_number": batchNumber}).Decode(&batch)

	if err != nil {
		return nil, err
	}

	return batch, nil

}
