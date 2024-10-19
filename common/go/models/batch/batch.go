package batch

import "go.mongodb.org/mongo-driver/bson/primitive"

type Batch struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string               `bson:"name" json:"name"`
	Tests []primitive.ObjectID `bson:"tests" json:"tests" ts_type:"string[]"`
}

func (batch *Batch) GetCollectionName() string {
	return "batches"
}
