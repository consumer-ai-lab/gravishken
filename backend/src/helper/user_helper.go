package helper

import (
	"context"
	"errors"
	"fmt"
	"server/src/types"
	"strconv"

	User "common/models/user"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserTestTime(Collection *mongo.Collection, Username string, TimeToIncrease int64) error {
	var user User.User

	err := Collection.FindOne(context.TODO(), bson.M{"name": Username}).Decode(&user)

	if err != nil {
		return err
	}

	userTest := user.Tests
	prevTimeElapsedUser := userTest.ElapsedTime
	userTest.ElapsedTime = prevTimeElapsedUser - 60*TimeToIncrease

	if userTest.ElapsedTime < 0 {
		userTest.ElapsedTime = 0
	}

	if userTest.ElapsedTime > 1797 {
		userTest.ElapsedTime = 1797
	}

	user.Tests = userTest

	Collection.ReplaceOne(context.TODO(), bson.M{"name": Username}, user)

	return nil
}

func UpdateBatchTestTime(Collection *mongo.Collection, Usernames []string, TimeToIncrease int64) error {
	for _, username := range Usernames {
		err := UpdateUserTestTime(Collection, username, TimeToIncrease)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateUserData(Collection *mongo.Collection, Model *User.UserUpdateRequest) error {

	var user User.User

	err := Collection.FindOne(context.TODO(), bson.M{"name": Model.Username}).Decode(&user)

	userTest := user.Tests
	if err != nil {
		return err
	}

	property := strings.ToLower(Model.Property)
	_ = property

	switch property {
	case "start_time":

		startTime, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.StartTime = startTime
		userTest.ElapsedTime = 0
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "reading_submission_received":
		userTest.ReadingSubmissionReceived = true
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "submission_received":
		userTest.SubmissionReceived = true
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "elapsed_time":
		elapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.ElapsedTime = elapsedTime.Unix()
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "reading_elapsed_time":
		readingElapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.ReadingElapsedTime = readingElapsedTime.Unix()
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "submission_folder_id":
		userTest.SubmissionFolderID = Model.Value[0]
		userTest.MergedFileID = Model.Value[1]
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "wpm":
		wpm, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.WPM = float64(wpm.Unix()) // Convert int64 to float64

		wmp_time, err := time.Parse(time.RFC3339, Model.Value[1])
		if err != nil {
			return err
		}
		userTest.WPMNormal = float64(wmp_time.Unix())

		wpm_normal, err := time.Parse(time.RFC3339, Model.Value[2])
		if err != nil {
			return err
		}
		userTest.WPMNormal = float64(wpm_normal.Unix())
		user.Tests = userTest
		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

	case "user_test_time":
		username := Model.Value[0]
		timeToIncrease, err := time.Parse(time.RFC3339, Model.Value[1])
		if err != nil {
			return err
		}
		err = UpdateUserTestTime(Collection, username, timeToIncrease.Unix()) // Convert timeToIncrease to int64
		if err != nil {
			return err
		}
	case "batch_test_time":
		batchNumber := Model.Value[0]
		timeToIncrease, err := time.Parse(time.RFC3339, Model.Value[1])
		if err != nil {
			return err
		}
		user, err := GetModelByBatchId(Collection, batchNumber, &User.User{})
		if err != nil {
			return err
		}

		var usernames []string
		for _, user := range user {
			usernames = append(usernames, user.(*User.User).Username)
		}

		err = UpdateBatchTestTime(Collection, usernames, timeToIncrease.Unix())

		if err != nil {
			return err
		}

		fmt.Printf("Batch id: %s, Time to increase: %d\n", batchNumber, timeToIncrease.Unix())
	default:
		return errors.New("invalid property")
	}

	return nil
}


func UpdateUser(collection *mongo.Collection, userRequest *User.UserModelUpdateRequest) error {

	var user User.User

	

	objectID, err := primitive.ObjectIDFromHex(userRequest.ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil{
		return err
	}

	user.Username = userRequest.Username
	user.Password = userRequest.Password
	user.TestPassword = userRequest.TestPassword
	user.Batch = userRequest.Batch

	collection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, user)

	return nil;

}


func GetBatchWiseList(Collection *mongo.Collection, BatchNumber string) ([]map[string]string, error) {
	var result []map[string]string
	user, err := GetModelByBatchId(Collection, BatchNumber, &User.User{})
	if err != nil {
		return nil, err
	}

	/*
				userData[i].username,
		        userData[i].merged_file_id,
		        userData[i].submission_folder_id,
	*/

	for _, user := range user {
		userMap := map[string]string{
			"username":             user.(*User.User).Username,
			"merged_file_id":       user.(*User.User).Tests.MergedFileID,
			"submission_folder_id": user.(*User.User).Tests.SubmissionFolderID,
		}
		result = append(result, userMap)
	}

	return result, nil

}

func GetBatchWiseListRoll(Collection *mongo.Collection, BatchNumber string, From, To int) ([]map[string]string, error) {
	var result []map[string]string
	user, err := GetModelByBatchId(Collection, BatchNumber, &User.User{})
	if err != nil {
		return nil, err
	}

	/*
		userData[i].username,
		userData[i].merged_file_id,
		userData[i].submission_folder_id,
		userData[i].resultDownloaded,
		userData[i].submission_received,
	*/

	for _, user := range user {
		username, _ := strconv.Atoi(user.(*User.User).Username) // Convert username to integer
		if username >= From && username <= To {
			userMap := map[string]string{
				"username":             user.(*User.User).Username,
				"merged_file_id":       user.(*User.User).Tests.MergedFileID,
				"submission_folder_id": user.(*User.User).Tests.SubmissionFolderID,
				"resultDownloaded":     strconv.FormatBool(user.(*User.User).Tests.ResultDownloaded),
				"submission_received":  strconv.FormatBool(user.(*User.User).Tests.SubmissionReceived),
			}
			result = append(result, userMap)
		}
	}

	return result, nil
}

func GetBatchDataForFrontend(Collection *mongo.Collection, BatchNumber string) ([]map[string]string, error) {
	var result []map[string]string
	user, err := GetModelByBatchId(Collection, BatchNumber, &User.User{})
	if err != nil {
		return nil, err
	}

	for _, user := range user {
		start_time := user.(*User.User).Tests.StartTime
		userArr := make(map[string]string)
		userArr["username"] = user.(*User.User).Username
		userArr["merged_file_id"] = user.(*User.User).Tests.MergedFileID
		userArr["submission_folder_id"] = user.(*User.User).Tests.SubmissionFolderID
		if start_time.IsZero() {
			userArr["status"] = "Present"
		} else {
			userArr["status"] = "Absent"
		}

		result = append(result, userArr)
	}

	return result, nil
}

func UserLogin(Collection *mongo.Collection, userRequest *User.UserLoginRequest) (string, error) {
	user, err := User.FindByUsername(Collection, userRequest.Username)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if user.Password != userRequest.Password {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userRequest.Username,
		"exp":      time.Now().Add(48 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("token"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetUserResultToDownloaded(Collection *mongo.Collection, request *User.UserBatchRequestData) error {
	user, err := Get_All_Models(Collection, &User.User{})
	if err != nil {
		return err
	}

	from := request.From
	to := request.To
	resultDownloaded := request.ResultDownloaded

	filered_users := []types.ModelInterface{}

	for _, user := range user {
		username, _ := strconv.Atoi(user.(*User.User).Username) // Convert username to integer
		if username >= from && username <= to {
			filered_users = append(filered_users, user)
		}
	}

	for _, filtered_user := range filered_users {
		if !filtered_user.(*User.User).Tests.SubmissionReceived {
			continue
		}

		filtered_user.(*User.User).Tests.ResultDownloaded = resultDownloaded
		err = Update_Model_By_ID(Collection, filtered_user.(*User.User).ID.Hex(), filtered_user)
		if err != nil {
			return err
		}

	}
	return nil
}

func ResetUserData(Collection *mongo.Collection, username string) error {
	user, err := User.FindByUsername(Collection, username)
	if err != nil {
		return err
	}

	/*
		userData.submission_received = false;
		userData.reading_submission_received = false;
		userData.reading_elapsed_time = 0;
		userData.elapsed_time = 0;
	*/

	user.Tests.SubmissionReceived = false
	user.Tests.ReadingSubmissionReceived = false
	user.Tests.ReadingElapsedTime = 0
	user.Tests.ElapsedTime = 0

	err = Update_Model_By_ID(Collection, user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return nil
}
