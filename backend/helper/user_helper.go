package helper

import (
	"context"
	"errors"
	"fmt"
	"gravtest/auth"
	"gravtest/types"
	"strconv"

	"gravtest/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func UpdateUserTestTime(Collection *mongo.Collection, Username string, TimeToIncrease int64) error {
	var user models.User

	err := Collection.FindOne(context.TODO(), bson.M{"username": Username}).Decode(&user)

	if err != nil {
		return err
	}

	userTest := user.Tests
	prevTimeElapsedUser := userTest.ElapsedTime
	userTest.ElapsedTime = prevTimeElapsedUser - 60* TimeToIncrease

	if userTest.ElapsedTime < 0 {
		userTest.ElapsedTime = 0
	}

	if userTest.ElapsedTime > 1797 {
		userTest.ElapsedTime = 1797
	}

	user.Tests = userTest

	Collection.ReplaceOne(context.TODO(), bson.M{"username": Username}, user)

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


func UpdateUserDate(Collection *mongo.Collection, Model *models.UserUpdateRequest) error {
	valid_request, err := auth.ValidRequestVerifier(Collection, Model.Token, Model.ApiKey)
	if err != nil {
		return err
	}

	if !valid_request {
		return errors.New("invalid request: Token or Apikey is invalid")
	}

	var userTest models.UserTest

	err = Collection.FindOne(context.TODO(), bson.M{"username": Model.Username}).Decode(&userTest)

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
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

	case "reading_submission_received":
		userTest.ReadingSubmissionReceived = true
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

	case "submission_received":
		userTest.SubmissionReceived = true
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

	case "elapsed_time":
		elapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.ElapsedTime = elapsedTime.Unix()
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

	case "reading_elapsed_time":
		readingElapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
		if err != nil {
			return err
		}
		userTest.ReadingElapsedTime = readingElapsedTime.Unix()
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

	case "submission_folder_id":
		userTest.SubmissionFolderID = Model.Value[0]
		userTest.MergedFileID = Model.Value[1]
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

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
		Collection.ReplaceOne(context.TODO(), bson.M{"username": Model.Username}, userTest)

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
		batchID := Model.Value[0]
		timeToIncrease, err := time.Parse(time.RFC3339, Model.Value[1])
		if err != nil {
			return err
		}
		user, err := GetModelByBatchId(Collection, batchID, &models.User{})
		if err != nil {
			return err
		}

		var usernames []string
		for _, user := range user {
			usernames = append(usernames, user.(*models.User).Username)	
		}

		err = UpdateBatchTestTime(Collection, usernames, timeToIncrease.Unix()) 

		if err != nil {
			return err
		}

		fmt.Printf("Batch id: %s, Time to increase: %d\n", batchID, timeToIncrease.Unix())	
	default:
		return errors.New("invalid property")
	}

	return nil
}


func GetBatchDataForFrontend(Collection *mongo.Collection, BatchID string) ([][]string, error){
	var result [][] string
	user, err := GetModelByBatchId(Collection, BatchID, &models.User{})
	if err != nil {
		return nil, err
	}

	for _, user := range user {
		start_time := user.(*models.User).Tests.StartTime
		userArr := []string{}
		if start_time.IsZero() {
			userArr = append(userArr, user.(*models.User).Username)
			userArr = append(userArr, user.(*models.User).Tests.MergedFileID)
			userArr = append(userArr, "Present")
			userArr = append(userArr, user.(*models.User).Tests.SubmissionFolderID)
		} else{
			userArr = append(userArr, user.(*models.User).Username)
			userArr = append(userArr, user.(*models.User).Tests.MergedFileID)
			userArr = append(userArr, "Absent")
			userArr = append(userArr, user.(*models.User).Tests.SubmissionFolderID)
		} 

		result = append(result, userArr)
	}

	return result, nil
}


// incomplete
func UserLogin(Collection *mongo.Collection, userRequest *models.UserLoginRequest) error{
	user, err := auth.FindUserByUsername(Collection, userRequest.Username)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	batch_data, err := GetQuestionPaper(Collection, userRequest.TestPassword)

	if err != nil {
		return err
	}

	if batch_data == nil{
		return errors.New("batch not found")
	}

	// Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(48 * time.Hour).Unix(),
    })

	tokenString, err := token.SignedString([]byte("token"))
    if err != nil {
        return err
    }

	user.Token = tokenString
	user.Password = userRequest.TestPassword  // I am still confused are there 2 passwords or only one password

	err = Update_Model_By_ID(Collection, user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return nil
}

type RequestData struct{
	from int
	to int
	resultDownloaded bool
}

func SetUserResultToDownloaded(Collection *mongo.Collection ,request *RequestData) error {
	user, err := Get_All_Models(Collection, &models.User{})
	if err != nil {
		return err
	}
	from := request.from
	to := request.to
	resultDownloaded := request.resultDownloaded
	filered_users := []types.ModelInterface{}
	for _, user := range user {
		username, _ := strconv.Atoi(user.(*models.User).Username) // Convert username to integer
		if username >= from && username <= to {
			filered_users = append(filered_users, user)
		}
	}

	for _, filtered_user := range filered_users{
		if !filtered_user.(*models.User).Tests.SubmissionReceived{
			continue
		}

		filtered_user.(*models.User).Tests.ResultDownloaded = resultDownloaded
		err = Update_Model_By_ID(Collection, filtered_user.(*models.User).ID.Hex(), filtered_user)
		if err != nil {
			return err
		}

	}
	return nil
}