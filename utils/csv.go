package utils

import (
	"encoding/csv"
	"os"
)

func CreateCSV(data [][]string, path string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()
	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return CannotWriteToFileError
		}
	}

	return nil
}
