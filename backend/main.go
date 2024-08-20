package main

import (
	"context"
	"fmt"
	"gravtest/helper"
	"gravtest/models"
	"gravtest/mongodb"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables
	dbName := os.Getenv("DATABASE")
	// batchCollectionName := os.Getenv("BATCH_COLLECTION")
	adminCollectionName := os.Getenv("ADMIN_COLLECTION")

	client, err := mongodb.Connect()
	if err != nil {
		log.Fatal("Error in connecting to MongoDB: ", err)
	}
	defer client.Disconnect(context.TODO())

	// BATCH_COLLECTION := client.Database(dbName).Collection(batchCollectionName)
	ADMIN_COLLECTION := client.Database(dbName).Collection(adminCollectionName)

	// utils.Add_CSVData_To_DB(BATCH_COLLECTION, "studentData.csv");
	// response, err := helper.Get_All_Batches(BATCH_COLLECTION)
	// if err != nil {
	// 	fmt.Println("Error in fetching all batches: ", err)
	// }

	// admin := models.Admin{
	// 	Username: "testing",
	// 	Password: "testing",
	// }

	// adminRequest := models.AdminRequest{
	// 	Username: "testing",
	// 	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RpbmciLCJleHAiOjE3MjQxMzkyMzN9.DXHhEcaCDfZiajTtjPcqNfXV6OVKPnKmy0NmcG1cJ2I",
	// }

	// err = helper.Delete_ALL_Model(ADMIN_COLLECTION)

	// if err != nil {
	// 	fmt.Println("Error in adding admin to database: ", err)
	// }

	// err = helper.AdminLogout(ADMIN_COLLECTION, &adminRequest)
	helper.ChangePassword(ADMIN_COLLECTION, &models.AdminChangePassword{
		Username:   "testing",
		NewPassword: "testing_new1",
	})

	// response, err := helper.Get_All_Models(ADMIN_COLLECTION, &models.Admin{})

	if err != nil {
		fmt.Println("Error in fetching all admins: ", err)
	}

}
