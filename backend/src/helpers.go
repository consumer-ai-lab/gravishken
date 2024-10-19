package main

import (
	"common"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	// "strconv"
	// "strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// func UpdateUserTestTime(Collection *mongo.Collection, Username string, TimeToIncrease int64) error {
// 	var user common.User

// 	err := Collection.FindOne(context.TODO(), bson.M{"name": Username}).Decode(&user)

// 	if err != nil {
// 		return err
// 	}

// 	userTest := user.Tests
// 	prevTimeElapsedUser := userTest.ElapsedTime
// 	userTest.ElapsedTime = prevTimeElapsedUser - 60*TimeToIncrease

// 	if userTest.ElapsedTime < 0 {
// 		userTest.ElapsedTime = 0
// 	}

// 	if userTest.ElapsedTime > 1797 {
// 		userTest.ElapsedTime = 1797
// 	}

// 	user.Tests = userTest

// 	Collection.ReplaceOne(context.TODO(), bson.M{"name": Username}, user)

// 	return nil
// }

// func UpdateBatchTestTime(Collection *mongo.Collection, Usernames []string, TimeToIncrease int64) error {
// 	for _, username := range Usernames {
// 		err := UpdateUserTestTime(Collection, username, TimeToIncrease)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func UpdateUserData(Collection *mongo.Collection, Model *common.UserUpdateRequest) error {

// 	var user common.User

// 	err := Collection.FindOne(context.TODO(), bson.M{"name": Model.Username}).Decode(&user)

// 	userTest := user.Tests
// 	if err != nil {
// 		return err
// 	}

// 	property := strings.ToLower(Model.Property)
// 	_ = property

// 	switch property {
// 	case "start_time":

// 		startTime, err := time.Parse(time.RFC3339, Model.Value[0])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.StartTime = startTime
// 		userTest.ElapsedTime = 0
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "reading_submission_received":
// 		userTest.ReadingSubmissionReceived = true
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "submission_received":
// 		userTest.SubmissionReceived = true
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "elapsed_time":
// 		elapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.ElapsedTime = elapsedTime.Unix()
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "reading_elapsed_time":
// 		readingElapsedTime, err := time.Parse(time.RFC3339, Model.Value[0])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.ReadingElapsedTime = readingElapsedTime.Unix()
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "submission_folder_id":
// 		userTest.SubmissionFolderID = Model.Value[0]
// 		userTest.MergedFileID = Model.Value[1]
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "wpm":
// 		wpm, err := time.Parse(time.RFC3339, Model.Value[0])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.WPM = float64(wpm.Unix()) // Convert int64 to float64

// 		wmp_time, err := time.Parse(time.RFC3339, Model.Value[1])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.WPMNormal = float64(wmp_time.Unix())

// 		wpm_normal, err := time.Parse(time.RFC3339, Model.Value[2])
// 		if err != nil {
// 			return err
// 		}
// 		userTest.WPMNormal = float64(wpm_normal.Unix())
// 		user.Tests = userTest
// 		Collection.ReplaceOne(context.TODO(), bson.M{"name": Model.Username}, user)

// 	case "user_test_time":
// 		username := Model.Value[0]
// 		timeToIncrease, err := time.Parse(time.RFC3339, Model.Value[1])
// 		if err != nil {
// 			return err
// 		}
// 		err = UpdateUserTestTime(Collection, username, timeToIncrease.Unix()) // Convert timeToIncrease to int64
// 		if err != nil {
// 			return err
// 		}
// 	case "batch_test_time":
// 		batchNumber := Model.Value[0]
// 		timeToIncrease, err := time.Parse(time.RFC3339, Model.Value[1])
// 		if err != nil {
// 			return err
// 		}
// 		user, err := GetModelByBatchId(Collection, batchNumber, &common.User{})
// 		if err != nil {
// 			return err
// 		}

// 		var usernames []string
// 		for _, user := range user {
// 			usernames = append(usernames, user.(*common.User).Username)
// 		}

// 		err = UpdateBatchTestTime(Collection, usernames, timeToIncrease.Unix())

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Printf("Batch id: %s, Time to increase: %d\n", batchNumber, timeToIncrease.Unix())
// 	default:
// 		return errors.New("invalid property")
// 	}

// 	return nil
// }

// func UpdateUser(collection *mongo.Collection, userRequest *common.UserModelUpdateRequest) error {

// 	var user common.User

// 	objectID, err := primitive.ObjectIDFromHex(userRequest.ID)
// 	if err != nil {
// 		return fmt.Errorf("invalid ID format: %v", err)
// 	}

// 	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
// 	if err != nil {
// 		return err
// 	}

// 	user.Username = userRequest.Username
// 	user.Password = userRequest.Password
// 	user.TestPassword = userRequest.TestPassword
// 	user.Batch = userRequest.Batch

// 	collection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, user)

// 	return nil

// }

// func GetBatchWiseList(Collection *mongo.Collection, BatchNumber string) ([]map[string]string, error) {
// 	var result []map[string]string
// 	user, err := GetModelByBatchId(Collection, BatchNumber, &common.User{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	/*
// 				userData[i].username,
// 		        userData[i].merged_file_id,
// 		        userData[i].submission_folder_id,
// 	*/

// 	for _, user := range user {
// 		userMap := map[string]string{
// 			"username":             user.(*common.User).Username,
// 			"merged_file_id":       user.(*common.User).Tests.MergedFileID,
// 			"submission_folder_id": user.(*common.User).Tests.SubmissionFolderID,
// 		}
// 		result = append(result, userMap)
// 	}

// 	return result, nil

// }

// func GetBatchWiseListRoll(Collection *mongo.Collection, BatchNumber string, From, To int) ([]map[string]string, error) {
// 	var result []map[string]string
// 	user, err := GetModelByBatchId(Collection, BatchNumber, &common.User{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	/*
// 		userData[i].username,
// 		userData[i].merged_file_id,
// 		userData[i].submission_folder_id,
// 		userData[i].resultDownloaded,
// 		userData[i].submission_received,
// 	*/

// 	for _, user := range user {
// 		username, _ := strconv.Atoi(user.(*common.User).Username) // Convert username to integer
// 		if username >= From && username <= To {
// 			userMap := map[string]string{
// 				"username":             user.(*common.User).Username,
// 				"merged_file_id":       user.(*common.User).Tests.MergedFileID,
// 				"submission_folder_id": user.(*common.User).Tests.SubmissionFolderID,
// 				"resultDownloaded":     strconv.FormatBool(user.(*common.User).Tests.ResultDownloaded),
// 				"submission_received":  strconv.FormatBool(user.(*common.User).Tests.SubmissionReceived),
// 			}
// 			result = append(result, userMap)
// 		}
// 	}

// 	return result, nil
// }

// func GetBatchDataForFrontend(Collection *mongo.Collection, BatchNumber string) ([]map[string]string, error) {
// 	var result []map[string]string
// 	user, err := GetModelByBatchId(Collection, BatchNumber, &common.User{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, user := range user {
// 		start_time := user.(*common.User).Tests.StartTime
// 		userArr := make(map[string]string)
// 		userArr["username"] = user.(*common.User).Username
// 		userArr["merged_file_id"] = user.(*common.User).Tests.MergedFileID
// 		userArr["submission_folder_id"] = user.(*common.User).Tests.SubmissionFolderID
// 		if start_time.IsZero() {
// 			userArr["status"] = "Present"
// 		} else {
// 			userArr["status"] = "Absent"
// 		}

// 		result = append(result, userArr)
// 	}

// 	return result, nil
// }

func UserLogin(Collection *mongo.Collection, userRequest *common.TUserLoginRequest) (string, error) {
	user, err := common.FindByUsername(Collection, userRequest.Username)

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

// func SetUserResultToDownloaded(Collection *mongo.Collection, request *common.UserBatchRequestData) error {
// 	user, err := Get_All_Models(Collection, &common.User{})
// 	if err != nil {
// 		return err
// 	}

// 	from := request.From
// 	to := request.To
// 	resultDownloaded := request.ResultDownloaded

// 	filered_users := []ModelInterface{}

// 	for _, user := range user {
// 		username, _ := strconv.Atoi(user.(*common.User).Username) // Convert username to integer
// 		if username >= from && username <= to {
// 			filered_users = append(filered_users, user)
// 		}
// 	}

// 	for _, filtered_user := range filered_users {
// 		if !filtered_user.(*common.User).Tests.SubmissionReceived {
// 			continue
// 		}

// 		filtered_user.(*common.User).Tests.ResultDownloaded = resultDownloaded
// 		err = Update_Model_By_ID(Collection, filtered_user.(*common.User).ID.Hex(), filtered_user)
// 		if err != nil {
// 			return err
// 		}

// 	}
// 	return nil
// }

// func ResetUserData(Collection *mongo.Collection, username string) error {
// 	user, err := common.FindByUsername(Collection, username)
// 	if err != nil {
// 		return err
// 	}

// 	/*
// 		userData.submission_received = false;
// 		userData.reading_submission_received = false;
// 		userData.reading_elapsed_time = 0;
// 		userData.elapsed_time = 0;
// 	*/

// 	user.Tests.SubmissionReceived = false
// 	user.Tests.ReadingSubmissionReceived = false
// 	user.Tests.ReadingElapsedTime = 0
// 	user.Tests.ElapsedTime = 0

// 	err = Update_Model_By_ID(Collection, user.ID.Hex(), user)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func GetTestsByBatch(batchCollection *mongo.Collection, testCollection *mongo.Collection, batchName string) ([]common.Test, error) {
	var batchDoc common.Batch
	err := batchCollection.FindOne(context.TODO(), bson.M{"name": batchName}).Decode(&batchDoc)
	if err != nil {
		return nil, fmt.Errorf("error finding batch: %v", err)
	}

	var tests []common.Test
	cursor, err := testCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": batchDoc.Tests}})
	if err != nil {
		return nil, fmt.Errorf("error finding tests: %v", err)
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &tests)
	if err != nil {
		return nil, fmt.Errorf("error decoding tests: %v", err)
	}

	return tests, nil
}

func GetTestByID(testCollection *mongo.Collection, testID primitive.ObjectID) (*common.Test, error) {
	var testDoc common.Test
	err := testCollection.FindOne(context.TODO(), bson.M{"_id": testID}).Decode(&testDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("error finding test: %v", err)
	}
	return &testDoc, nil
}

func GetBatchByBatchNumber(Collection *mongo.Collection, batchNumber string) (ModelInterface, error) {
	var batch ModelInterface

	err := Collection.FindOne(context.TODO(), bson.M{"batch_number": batchNumber}).Decode(&batch)

	if err != nil {
		return nil, err
	}

	return batch, nil

}

func Add_Model_To_DB(Collection *mongo.Collection, Model ModelInterface) error {
	_, err := Collection.InsertOne(context.TODO(), Model)

	if err != nil {
		fmt.Println("Error in adding Model to database: ", err)
		return err
	}

	return nil
}

func Get_All_Models(collection *mongo.Collection, modelType ModelInterface) ([]ModelInterface, error) {
	fmt.Println("Fetching all models from database...")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error in fetching all models: ", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []ModelInterface

	// Create a new slice of the concrete type
	sliceType := reflect.SliceOf(reflect.TypeOf(modelType).Elem())
	concreteSlice := reflect.MakeSlice(sliceType, 0, 0)
	concreteSlicePtr := reflect.New(concreteSlice.Type())
	concreteSlicePtr.Elem().Set(concreteSlice)

	// Decode into the concrete slice
	if err := cursor.All(context.TODO(), concreteSlicePtr.Interface()); err != nil {
		fmt.Println("Error in decoding models: ", err)
		return nil, err
	}

	// Convert concrete slice to []ModelInterface
	concreteSliceVal := concreteSlicePtr.Elem()
	for i := 0; i < concreteSliceVal.Len(); i++ {
		item := concreteSliceVal.Index(i).Addr().Interface().(ModelInterface)
		results = append(results, item)
	}

	return results, nil
}

func GetModelById(collection *mongo.Collection, ID string, modelType ModelInterface) (ModelInterface, error) {
	fmt.Println("Fetching model from database...")
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var result ModelInterface
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		fmt.Println("Error in fetching model: ", err)
		return nil, err
	}

	fmt.Println("Model in DB:", result)
	return result, nil
}

func GetModelByBatchId(collection *mongo.Collection, batchNumber string, modelType ModelInterface) ([]ModelInterface, error) {
	fmt.Println("Fetching model from database...")

	var results []ModelInterface
	cursor, err := collection.Find(context.TODO(), bson.M{"batch": batchNumber})
	if err != nil {
		fmt.Println("Error in fetching model: ", err)
		return nil, err
	}

	sliceType := reflect.SliceOf(reflect.TypeOf(modelType).Elem())
	concreteSlice := reflect.MakeSlice(sliceType, 0, 0)
	concreteSlicePtr := reflect.New(concreteSlice.Type())
	concreteSlicePtr.Elem().Set(concreteSlice)

	// Decode into the concrete slice
	if err := cursor.All(context.TODO(), concreteSlicePtr.Interface()); err != nil {
		fmt.Println("Error in decoding models: ", err)
		return nil, err
	}

	// Convert concrete slice to []ModelInterface
	concreteSliceVal := concreteSlicePtr.Elem()
	for i := 0; i < concreteSliceVal.Len(); i++ {
		item := concreteSliceVal.Index(i).Addr().Interface().(ModelInterface)
		results = append(results, item)
	}

	return results, nil
}

func Delete_Model_By_ID(Collection *mongo.Collection, ID string) error {

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	result, err := Collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error in deleting Model from database: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with ID: %s", ID)
	}

	log.Printf("Model with ID %s deleted successfully", ID)
	return nil
}

func Delete_ALL_Model(Collection *mongo.Collection) error {

	_, err := Collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error in deleting Model from database: ", err)
		return err
	}

	log.Default().Println("Model deleted successfully")

	return nil
}

func Update_Model_By_ID(Collection *mongo.Collection, ID string, Model ModelInterface) error {

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	fmt.Println("Updating Model with ID: ", ID)

	result, err := Collection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, Model)
	if err != nil {
		log.Printf("Error in updating Model from database: %v", err)
		return err
	}

	fmt.Println("Model updated successfully")

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID: %s", ID)
	}

	log.Printf("Model with ID %s updated successfully", ID)
	return nil
}

func RegisterAdmin(Collection *mongo.Collection, Admin ModelInterface) error {

	password := Admin.(*common.Admin).Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}

	Admin.(*common.Admin).Password = string(hashedPassword)

	Add_Model_To_DB(Collection, Admin)
	return nil
}

func AdminLogin(Collection *mongo.Collection, Admin ModelInterface) (string, error) {
	username := Admin.(*common.Admin).Username
	password := Admin.(*common.Admin).Password
	secretKey := []byte("TODO:add-a-secret-key-from-env")

	var user common.Admin
	err := Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("admin not found")
		}
		return "", fmt.Errorf("error finding admin: %v", err)
	}

	log.Default().Printf("Provided username: %s and password: %s\nDatabase usernam: %s and password: %s", username, password, user.Username, user.Password)

	// Compare the hashed password with the plaintext password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	expirationTime := time.Now().Add(48 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return tokenString, nil
}

func ValidateAdminToken(tokenString string) (*Claims, error) {
	secretKey := []byte("TODO:add-a-secret-key-from-env") // Use the same secret key as in AdminLogin

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func UpdateTypingTestText(collection *mongo.Collection, testID string, typingText string) error {
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": testID, "type": "typing"},
		bson.M{"$set": bson.M{"typingText": typingText}},
	)
	if err != nil {
		return fmt.Errorf("error updating typing test text: %v", err)
	}

	return nil
}
