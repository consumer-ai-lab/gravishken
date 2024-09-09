package test

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType string

const (
	PDF FileType = "pdf"
	DOC FileType = "doc"
	TXT FileType = "txt"
)

type Test struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileType       FileType           `bson:"fileType" json:"fileType" binding:"required"`
	TimeSlot       time.Time          `bson:"timeSlot" json:"timeSlot" binding:"required"`
	Password       string             `bson:"password" json:"password" binding:"required"`
	DriveID        string             `bson:"driveId,omitempty" json:"driveId,omitempty"`
	TypingTestText string             `bson:"typingTestText" json:"typingTestText"`
	BatchNumber    string             `bson:"batch" json:"batch" binding:"required"`
}

func (test *Test) GetCollectionName() string {
	return "tests"
}
