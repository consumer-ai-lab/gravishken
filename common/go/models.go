package common

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (test *Test) GetCollectionName() string {
	return "tests"
}

func (user *User) GetCollectionName() string {
	return "users"
}

func (userTest *UserSubmission) GetCollectionName() string {
	return "submissions"
}

func (admin *Admin) GetCollectionName() string {
	return "admins"
}

func (batch *Batch) GetCollectionName() string {
	return "batches"
}

// primitive id converted to string
type ID = string

type Admin struct {
	Id       ID
	Username string
	Password string
}

type AdminRequest struct {
	Username string
	Token    string
}

type Batch struct {
	Id    ID
	Name  string
	Tests []ID
}

type Test struct {
	Id         ID
	Type       TestType
	Duration   int
	File       string
	TypingText string
}

type User struct {
	Id       ID
	Username string
	// TODO: plaintext password yo!
	// passwords should be stored in another table hashed
	Password  string
	BatchName string
}

type UserSubmission struct {
	UserId ID
	TestId ID

	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime int64

	WPM       float64
	WPMNormal float64

	//?
	ReadingSubmissionReceived bool
	ReadingElapsedTime        int64
	SubmissionReceived        bool
	ResultDownloaded          bool
	MergedFileID              string
	SubmissionFolderID        string
}

type UserModelUpdateRequest struct {
	Id           ID
	Username     string
	Password     string
	TestPassword string
	Batch        string
}

type UserBatchRequestData struct {
	From             int
	To               int
	ResultDownloaded bool
}

type UserLoginResponse struct {
	Jwt  string
	User User
}

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

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return client.Database(dbName).Collection(collectionName)
}
