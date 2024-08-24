package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func Read_CSV(file_path string) ([]map[string]string, map[string]bool) {
    file, err := os.Open(file_path)
    if err != nil {
        log.Fatalf("failed to open file: %v", err)
    }
    defer file.Close()

	reader := csv.NewReader(file)

	csv_data := make([]map[string]string, 0)
	unique_batches := make(map[string]bool, 0)

	for {
		csv_single_data := make(map[string]string, 0)
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to read record: %v", err)
		}
		lastColumn := record[len(record)-1]
		unique_batches[lastColumn] = true
		csv_single_data["sr_no"] = record[0]
		csv_single_data["roll_no"] = record[1]
		csv_single_data["name"] = record[2]
		csv_single_data["father_name"] = record[3]
		csv_single_data["neis_no"] = record[4]
		csv_single_data["designation"] = record[5]
		csv_single_data["area"] = record[6]
		csv_single_data["slot"] = record[7]

		csv_data = append(csv_data, csv_single_data)
	}
	delete(unique_batches, "slot_allocated")

	return csv_data[1:], unique_batches
}