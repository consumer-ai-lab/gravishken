package utils


import (
    "encoding/csv"
    "log"
    "os"
	"io"
)

func Read_CSV(file_path string) [][]string{
    file, err := os.Open(file_path)
    if err != nil {
        log.Fatalf("failed to open file: %v", err)
    }
    defer file.Close()

	reader := csv.NewReader(file)

	csv_data := make([][]string, 0)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to read record: %v", err)
		}
		csv_data = append(csv_data, record)
	}

	return csv_data
}