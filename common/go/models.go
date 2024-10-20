package common

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (test *Test) GetCollectionName() string {
	return "tests"
}

func (user *User) GetCollectionName() string {
	return "users"
}

func (userTest *TestSubmission) GetCollectionName() string {
	return "submissions"
}

func (admin *Admin) GetCollectionName() string {
	return "admins"
}

func (batch *Batch) GetCollectionName() string {
	return "batches"
}

// primitive id converted to string
// type ID = string
type ID = primitive.ObjectID

type Admin struct {
	Id       ID `bson:"_id,omitempty" ts_type:"string"`
	Username string
	Password string
}

// type AdminRequest struct {
// 	Username string
// 	Token    string
// }

type Batch struct {
	Id    ID `bson:"_id,omitempty" ts_type:"string"`
	Name  string
	Tests []ID `ts_type:"string[]"`
}

type MCQ struct {
	Question string
	Options  []string
	Answer   string
}
type Test struct {
	Id       ID `bson:"_id,omitempty" ts_type:"string"`
	TestName string
	Duration int

	Type       TestType
	FilePath   string `bson:"file,omitempty" json:"FilePath,omitempty"`
	TypingText string `bson:"typingtext,omitempty" json:"TypingText,omitempty"`
	McqJson    string `bson:"mcqjson,omitempty" json:"McqJson,omitempty"`
}

type User struct {
	Id       ID `bson:"_id,omitempty" ts_type:"string"`
	Username string
	// TODO: plaintext password yo!
	// passwords should be stored in another table hashed
	Password string
	Batch    string
}

type AppTestInfo struct {
	FileData string
}
type McqTestInfo struct {
	Data string
}
type TypingTestInfo struct {
	WPM float64
}
type TestInfo struct {
	Type           TestType
	TypingTestInfo *TypingTestInfo `bson:"typingtestinfo,omitempty" json:"TypingTestInfo,omitempty"`
	McqTestInfo    *McqTestInfo    `bson:"mcqtestinfo,omitempty" json:"McqTestInfo,omitempty"`
	DocxTestInfo   *AppTestInfo    `bson:"docxtestinfo,omitempty" json:"DocxTestInfo,omitempty"`
	ExcelTestInfo  *AppTestInfo    `bson:"exceltestinfo,omitempty" json:"ExcelTestInfo,omitempty"`
	PptTestInfo    *AppTestInfo    `bson:"ppttestinfo,omitempty" json:"PptTestInfo,omitempty"`
}

type TestSubmission struct {
	UserId ID `ts_type:"string"`
	TestId ID `ts_type:"string"`

	StartTime time.Time
	EndTime   time.Time

	TestInfo TestInfo
}

// type UserModelUpdateRequest struct {
// 	Id           ID
// 	Username     string
// 	Password     string
// 	TestPassword string
// 	Batch        string
// }

// type UserBatchRequestData struct {
// 	From             int
// 	To               int
// 	ResultDownloaded bool
// }

type UserLoginResponse struct {
	Jwt  string
	User User
}

type TestType string

const (
	TypingTest TestType = "typing"
	DocxTest   TestType = "docx"
	ExcelTest  TestType = "xlsx"
	PptTest    TestType = "pptx"
	MCQTest    TestType = "mcq"
)

func (self TestType) TSName() string {
	switch self {
	case TypingTest:
		return "TypingTest"
	case DocxTest:
		return "DocxTest"
	case ExcelTest:
		return "ExcelTest"
	case PptTest:
		return "PptTest"
	case MCQTest:
		return "MCQTest"
	default:
		return "Unknown"
	}
}

func FindAdminByUsername(collection *mongo.Collection, username string) (*Admin, error) {
	filter := bson.M{"username": username}

	var admin Admin
	err := collection.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func FindByUsername(Collection *mongo.Collection, userName string) (*User, error) {

	filter := bson.M{"username": userName}

	var user User
	err := Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *Test) SetMCQQuestions(questions []MCQ) error {
	jsonData, err := json.Marshal(questions)
	if err != nil {
		return err
	}
	t.McqJson = string(jsonData)
	return nil
}

func (t *Test) GetMCQQuestions() ([]MCQ, error) {
	if t.McqJson == "" {
		return nil, nil
	}
	var questions []MCQ
	err := json.Unmarshal([]byte(t.McqJson), &questions)
	return questions, err
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return client.Database(dbName).Collection(collectionName)
}
