package test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestType string

const (
	TypingTest TestType = "typing"
	DocxTest   TestType = "docx"
	ExcelTest  TestType = "excel"
	WordTest   TestType = "word"
)

func (self TestType) TSName() string {
	switch self {
	case TypingTest:
		return "TypingTest"
	case DocxTest:
		return "DocxTest"
	case ExcelTest:
		return "ExcelTest"
	case WordTest:
		return "WordTest"
	default:
		return "Unknown"
	}
}

type Test struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type       TestType           `bson:"type" json:"type"`
	Duration   int                `bson:"duration" json:"duration"`
	File       string             `bson:"file,omitempty" json:"file,omitempty"`
	TypingText string             `bson:"typingText,omitempty" json:"typingText,omitempty"`
}

func (test *Test) GetCollectionName() string {
	return "tests"
}
