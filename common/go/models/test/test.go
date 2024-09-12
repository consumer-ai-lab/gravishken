package test

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType string

const (
	PPTX FileType = "pptx"
	DOCX FileType = "docx"
	XLSX FileType = "xlsx"
)

func (self FileType) TSName() string {
	switch self {
	case DOCX:
		return "DOCX"
	case XLSX:
		return "XLSX"
	case PPTX:
		return "PPTX"
	default:
		return "Unknown"
	}
}

type TestType string

const (
	TypingTest TestType = "typing"
	FileTest   TestType = "file"
)

func (self TestType) TSName() string {
	switch self {
	case TypingTest:
		return "TypingTest"
	case FileTest:
		return "FileTest"
	default:
		return "Unknown"
	}
}

type BatchTests struct {
	BatchId      primitive.ObjectID `bson:"batchId" json:"batchId"`
	Tests        []Test             `bson:"tests" json:"tests"`
	TestDuration int                `bson:"testDuration" json:"testDuration"`
	Password     string             `bson:"password" json:"password" binding:"required"`
	StartTime    time.Time          `bson:"startTime" json:"startTime"`
	EndTime      time.Time          `bson:"endTime" json:"endTime"`
}

type Test struct {
	TestId   primitive.ObjectID `bson:"testId" json:"testId"`
	TestType TestType           `bson:"testType" json:"testType"`

	// application tests
	FileType FileType `bson:"fileType,omitempty" json:"fileType,omitempty"`

	// typing test
	TypingTestText string `bson:"typingTestText,omitempty" json:"typingTestText,omitempty"`
}

func (test *BatchTests) GetCollectionName() string {
	return "tests"
}
