package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

func (admin *Admin) GetCollectionName() string {
	return "admins"
}

type AdminRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
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

type Batch struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string               `bson:"name" json:"name"`
	Tests []primitive.ObjectID `bson:"tests" json:"tests" ts_type:"string[]"`
}

func (batch *Batch) GetCollectionName() string {
	return "batches"
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

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required"`
	BatchName string             `bson:"batch_name" json:"batch_name" binding:"required"`
}

type UserSubmission struct {
	UserID                    primitive.ObjectID `bson:"user_id" json:"user_id"`
	TestID                    primitive.ObjectID `bson:"test" json:"test"`
	StartTime                 time.Time          `bson:"startTime" json:"startTime"`
	EndTime                   time.Time          `bson:"endTime" json:"endTime"`
	ElapsedTime               int64              `bson:"elapsedTime" json:"elapsedTime"` // Stored in seconds
	SubmissionReceived        bool               `bson:"submissionReceived" json:"submissionReceived"`
	ReadingElapsedTime        int64              `bson:"readingElapsedTime" json:"readingElapsedTime"` // Stored in seconds
	ReadingSubmissionReceived bool               `bson:"readingSubmissionReceived" json:"readingSubmissionReceived"`
	SubmissionFolderID        string             `bson:"submissionFolderId" json:"submissionFolderId"`
	MergedFileID              string             `bson:"mergedFileId" json:"mergedFileId"`
	WPM                       float64            `bson:"wpm" json:"wpm"`
	WPMNormal                 float64            `bson:"wpmNormal" json:"wpmNormal"`
	ResultDownloaded          bool               `bson:"resultDownloaded" json:"resultDownloaded"`
}

type UserModelUpdateRequest struct {
	ID           string `bson:"id" json:"id"`
	Username     string `bson:"username" json:"username" binding:"required"`
	Password     string `bson:"password" json:"password" binding:"required"`
	TestPassword string `bson:"testPassword" json:"testPassword" binding:"required"`
	Batch        string `bson:"batch" json:"batch" binding:"required"`
}

func (user *User) GetCollectionName() string {
	return "users"
}

func (userTest *UserSubmission) GetCollectionName() string {
	return "user_tests"
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

type UserBatchRequestData struct {
	From             int
	To               int
	ResultDownloaded bool
}

type UserUpdateRequest struct {
	Username string   `json:"username"`
	Property string   `json:"property"`
	Value    []string `json:"value"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return client.Database(dbName).Collection(collectionName)
}
