package test

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestType string

const (
	TypingTest TestType = "typing"
	DocxTest   TestType = "docx"
	ExcelTest  TestType = "excel"
	WordTest   TestType = "word"
	MCQTest    TestType = "mcq"
)

// MCQ represents the structure of a single multiple choice question
type MCQ struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

type Test struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TestName   string             `bson:"testName" json:"testName"`
	Type       TestType          `bson:"type" json:"type"`
	Duration   int               `bson:"duration" json:"duration"`
	FilePath   string           `bson:"file,omitempty" json:"file,omitempty"`
	TypingText string           `bson:"typingText,omitempty" json:"typingText,omitempty"`
	MCQJSON    string           `bson:"mcqJSON,omitempty" json:"mcqJSON,omitempty"`
}

// SetMCQQuestions sets the MCQ questions by marshaling them to JSON
func (t *Test) SetMCQQuestions(questions []MCQ) error {
	jsonData, err := json.Marshal(questions)
	if err != nil {
		return err
	}
	t.MCQJSON = string(jsonData)
	return nil
}

// GetMCQQuestions retrieves the MCQ questions by unmarshaling from JSON
func (t *Test) GetMCQQuestions() ([]MCQ, error) {
	if t.MCQJSON == "" {
		return nil, nil
	}
	var questions []MCQ
	err := json.Unmarshal([]byte(t.MCQJSON), &questions)
	return questions, err
}

func (t TestType) TSName() string {
	switch t {
	case TypingTest:
		return "TypingTest"
	case DocxTest:
		return "DocxTest"
	case ExcelTest:
		return "ExcelTest"
	case WordTest:
		return "WordTest"
	case MCQTest:
		return "MCQTest"
	default:
		return "Unknown"
	}
}

func (test *Test) GetCollectionName() string {
	return "tests"
}