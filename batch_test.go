package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCurlBatch(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "curl_batch_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create curl template file
	curlFile := filepath.Join(tmpDir, "curl.txt")
	curlContent := `curl -X POST -H "Content-Type: application/json" -d '{"name": "${NAME}"}' https://api.example.com`
	err = os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create curl file: %v", err)
	}

	// Create CSV file
	csvFile := filepath.Join(tmpDir, "data.csv")
	csvContent := `NAME,EMAIL
田中太郎,tanaka@example.com
佐藤花子,sato@example.com`
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}

	// Create output file path
	outputFile := filepath.Join(tmpDir, "output.txt")

	// Test NewCurlBatch
	batch, err := NewCurlBatch(curlFile, csvFile, outputFile, 1000)
	if err != nil {
		t.Fatalf("NewCurlBatch failed: %v", err)
	}
	defer batch.OutputFile.Close()

	if batch.CurlTemplate != curlContent {
		t.Errorf("Expected curl template %q, got %q", curlContent, batch.CurlTemplate)
	}

	if len(batch.CSVData) != 2 {
		t.Errorf("Expected 2 CSV records, got %d", len(batch.CSVData))
	}

	if batch.SleepMsec != 1000 {
		t.Errorf("Expected sleep 1000, got %d", batch.SleepMsec)
	}

	// Verify output file was created for append
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file should be created")
	}
}

func TestNewCurlBatchWithNonExistentFiles(t *testing.T) {
	tests := []struct {
		name        string
		curlFile    string
		csvFile     string
		outputFile  string
		expectError bool
	}{
		{
			name:        "Non-existent curl file",
			curlFile:    "nonexistent.txt",
			csvFile:     "sample/hogehoge.csv",
			outputFile:  "output.txt",
			expectError: true,
		},
		{
			name:        "Non-existent CSV file",
			curlFile:    "sample/curl.txt",
			csvFile:     "nonexistent.csv",
			outputFile:  "output.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCurlBatch(tt.curlFile, tt.csvFile, tt.outputFile, 0)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestNewCurlBatchOutputFileAppend(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "curl_batch_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create curl template file
	curlFile := filepath.Join(tmpDir, "curl.txt")
	curlContent := `curl -X GET https://api.example.com`
	err = os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create curl file: %v", err)
	}

	// Create CSV file
	csvFile := filepath.Join(tmpDir, "data.csv")
	csvContent := `NAME
test`
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}

	// Create existing output file
	outputFile := filepath.Join(tmpDir, "output.txt")
	existingContent := "Existing content\n"
	err = os.WriteFile(outputFile, []byte(existingContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing output file: %v", err)
	}

	// Test NewCurlBatch with existing output file
	batch, err := NewCurlBatch(curlFile, csvFile, outputFile, 0)
	if err != nil {
		t.Fatalf("NewCurlBatch failed: %v", err)
	}

	// Write some content to verify append mode
	_, err = batch.OutputFile.WriteString("New content\n")
	if err != nil {
		t.Fatalf("Failed to write to output file: %v", err)
	}
	batch.OutputFile.Close()

	// Read the file and verify content was appended
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "Existing content\nNew content\n"
	if string(content) != expectedContent {
		t.Errorf("Expected %q, got %q", expectedContent, string(content))
	}
}

func TestNewCurlBatchDefaultSleep(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "curl_batch_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create minimal files
	curlFile := filepath.Join(tmpDir, "curl.txt")
	err = os.WriteFile(curlFile, []byte("curl https://api.example.com"), 0644)
	if err != nil {
		t.Fatalf("Failed to create curl file: %v", err)
	}

	csvFile := filepath.Join(tmpDir, "data.csv")
	err = os.WriteFile(csvFile, []byte("NAME\ntest"), 0644)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.txt")

	// Test with sleep value 0
	batch, err := NewCurlBatch(curlFile, csvFile, outputFile, 0)
	if err != nil {
		t.Fatalf("NewCurlBatch failed: %v", err)
	}
	defer batch.OutputFile.Close()

	if batch.SleepMsec != 0 {
		t.Errorf("Expected sleep 0, got %d", batch.SleepMsec)
	}
}
