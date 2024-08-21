package main

import "time"

// added by kurve, just for testing
type UserTest struct {
	UserID                    string    `bson:"userId" json:"userId"`
	TestID                    string    `bson:"test" json:"test"`
	StartTime                 time.Time `bson:"startTime" json:"startTime"`
	EndTime                   time.Time `bson:"endTime" json:"endTime"`
	ElapsedTime               int64     `bson:"elapsedTime" json:"elapsedTime"` // Stored in seconds
	SubmissionReceived        bool      `bson:"submissionReceived" json:"submissionReceived"`
	ReadingElapsedTime        int64     `bson:"readingElapsedTime" json:"readingElapsedTime"` // Stored in seconds
	ReadingSubmissionReceived bool      `bson:"readingSubmissionReceived" json:"readingSubmissionReceived"`
	SubmissionFolderID        string    `bson:"submissionFolderId" json:"submissionFolderId"`
	MergedFileID              string    `bson:"mergedFileId" json:"mergedFileId"`
	WPM                       float64   `bson:"wpm" json:"wpm"`
	WPMNormal                 float64   `bson:"wpmNormal" json:"wpmNormal"`
	ResultDownloaded          bool      `bson:"resultDownloaded" json:"resultDownloaded"`
}
