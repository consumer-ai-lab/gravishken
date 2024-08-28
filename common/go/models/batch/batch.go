package batch

import "go.mongodb.org/mongo-driver/bson/primitive"

type Batch struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"batchName" json:"batchName"`
}

func (batch *Batch) GetCollectionName() string {
	return "batch"
}
