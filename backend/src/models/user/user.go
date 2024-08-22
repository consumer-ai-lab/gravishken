package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Username     string             `bson:"username" json:"username" binding:"required"`
	Password     string             `bson:"password" json:"password" binding:"required"`
	TestPassword string             `bson:"testPassword" json:"testPassword" binding:"required"`
	Token        string             `bson:"token" json:"token"`
	BatchID      primitive.ObjectID `bson:"batch" json:"batch,omitempty" binding:"required"`
	Tests        UserTest           `bson:"tests" json:"tests,omitempty"`
}

type UserTest struct {
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

func (user *User) GetCollectionName() string {
	return "users"
}

func (userTest *UserTest) GetCollectionName() string {
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

type UserUpdateRequest struct {
	Username string   `json:"username"`
	Token    string   `json:"token"`
	ApiKey   string   `json:"apiKey"`
	Property string   `json:"property"`
	Value    []string `json:"value"`
}

type UserLoginRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	TestPassword string `json:"testPassword"`
}
