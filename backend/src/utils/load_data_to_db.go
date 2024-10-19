package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func Add_CSVData_To_DB(collection *mongo.Collection, file_path string) error {

	// csvData := Read_CSV(file_path)
	// fmt.Println(csvData)

	// ****** Change logic here, add users to db and not batches ******

	// batches := make(map[string]*batch.Batch)

	// for _, row := range csvData[1:] { // Skip header row
	// 	lastCol := row[len(row)-1]
	// 	splitNumber := strings.Fields(lastCol)
	// 	if len(splitNumber) < 2 {
	// 		continue
	// 	}
	// 	batchNumber := splitNumber[1]
	// 	if _, exists := batches[batchNumber]; !exists {
	// 		batches[batchNumber] = &batch.Batch{Name: batchNumber}
	// 	}
	// }

	// for key := range batches {
	// 	err := helper.Add_Model_To_DB(batch_client, batches[key])
	// 	if err != nil {
	// 		return err
	// 	}

	// }

	return nil

}
