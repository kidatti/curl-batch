package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// readCSVData reads a CSV file and returns data as a slice of maps
// where each map represents a row with column headers as keys
func readCSVData(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	headers := records[0]
	var data []map[string]string

	for i, record := range records[1:] {
		if len(record) != len(headers) {
			return nil, fmt.Errorf("record %d has %d fields, expected %d", i+2, len(record), len(headers))
		}

		row := make(map[string]string)
		for j, value := range record {
			row[headers[j]] = value
		}
		data = append(data, row)
	}

	return data, nil
}
