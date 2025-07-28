package main

import (
	"fmt"
	"os"
	"time"
)

// CurlBatch represents a batch of curl requests to be executed
type CurlBatch struct {
	CurlTemplate string
	CSVData      []map[string]string
	OutputFile   *os.File
	SleepMsec    int
}

// NewCurlBatch creates a new CurlBatch instance
func NewCurlBatch(curlFile, csvFile, outputFile string, sleepMsec int) (*CurlBatch, error) {
	curlTemplate, err := readCurlTemplate(curlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read curl template: %w", err)
	}

	csvData, err := readCSVData(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	output, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}

	return &CurlBatch{
		CurlTemplate: curlTemplate,
		CSVData:      csvData,
		OutputFile:   output,
		SleepMsec:    sleepMsec,
	}, nil
}

// Run executes all curl requests in the batch
func (cb *CurlBatch) Run() error {
	defer cb.OutputFile.Close()

	for i, row := range cb.CSVData {
		curlCommand := cb.replaceTemplate(cb.CurlTemplate, row)

		fmt.Fprintf(cb.OutputFile, "=== Request %d ===\n", i+1)
		fmt.Fprintf(cb.OutputFile, "Command: %s\n", curlCommand)
		fmt.Fprintf(cb.OutputFile, "Data: %+v\n", row)

		result, err := cb.executeRequest(curlCommand)
		if err != nil {
			fmt.Fprintf(cb.OutputFile, "Error: %s\n", err)
		} else {
			fmt.Fprintf(cb.OutputFile, "Result:\n%s\n", result)
		}

		fmt.Fprintf(cb.OutputFile, "\n")

		fmt.Printf("Completed request %d/%d\n", i+1, len(cb.CSVData))

		// Sleep between requests if specified
		if cb.SleepMsec > 0 && i < len(cb.CSVData)-1 {
			time.Sleep(time.Duration(cb.SleepMsec) * time.Millisecond)
		}
	}

	return nil
}
