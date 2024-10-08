package utils

import (
	"encoding/base64"
	"encoding/csv"
	"io"
	"os"
)

var ParseCSVFile = parseCSVFile

func parseCSVFile(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var records [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				// reached the end of the file
				break
			}

			// any other type of error must finish the execution
			return nil, err
		}

		records = append(records, record)
	}
	
	return records[1:], nil
}

func EncodeFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file contents
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Encode the file content to base64
	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	// Print the base64 encoded string
	return encoded, nil
}
