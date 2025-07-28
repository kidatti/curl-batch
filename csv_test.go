package main

import (
	"os"
	"reflect"
	"testing"
)

func TestReadCSVData(t *testing.T) {
	// Create a temporary CSV file
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	csvContent := `NAME,EMAIL,AGE
田中太郎,tanaka@example.com,30
佐藤花子,sato@example.com,25`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{"NAME": "田中太郎", "EMAIL": "tanaka@example.com", "AGE": "30"},
		{"NAME": "佐藤花子", "EMAIL": "sato@example.com", "AGE": "25"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataEmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Close()

	_, err = readCSVData(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for empty CSV file, but got none")
	}
}

func TestReadCSVDataOnlyHeaders(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("NAME,EMAIL,AGE")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result for headers-only CSV, got %v", result)
	}
}

func TestReadCSVDataMismatchedColumns(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	csvContent := `NAME,EMAIL,AGE
田中太郎,tanaka@example.com,30,extra_field`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = readCSVData(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for mismatched columns, but got none")
	}
}

func TestReadCSVDataWithSpaces(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	csvContent := `NAME,EMAIL,AGE
 田中太郎 , tanaka@example.com , 30 
佐藤花子,sato@example.com,25`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	// Spaces should be preserved as-is
	expected := []map[string]string{
		{"NAME": " 田中太郎 ", "EMAIL": " tanaka@example.com ", "AGE": " 30 "},
		{"NAME": "佐藤花子", "EMAIL": "sato@example.com", "AGE": "25"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataFileNotExists(t *testing.T) {
	_, err := readCSVData("nonexistent_file.csv")
	if err == nil {
		t.Error("Expected error for non-existent file, but got none")
	}
}

func TestReadCSVDataSpecialCharacters(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	csvContent := `NAME,MESSAGE,SPECIAL
"田中太郎","Hello, World!","Special chars: @#$%^&*()"
"佐藤花子","こんにちは","Unicode: 🌸🗾"`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{"NAME": "田中太郎", "MESSAGE": "Hello, World!", "SPECIAL": "Special chars: @#$%^&*()"},
		{"NAME": "佐藤花子", "MESSAGE": "こんにちは", "SPECIAL": "Unicode: 🌸🗾"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataEscapedQuotes(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV with escaped double quotes using double quote escaping
	csvContent := `NAME,MESSAGE,DESCRIPTION
"John ""Johnny"" Doe","He said ""Hello World""","This is a ""test"" message"
"Jane Smith","She replied ""Hi there!""","Another ""quoted"" example"`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":        `John "Johnny" Doe`,
			"MESSAGE":     `He said "Hello World"`,
			"DESCRIPTION": `This is a "test" message`,
		},
		{
			"NAME":        "Jane Smith",
			"MESSAGE":     `She replied "Hi there!"`,
			"DESCRIPTION": `Another "quoted" example`,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataCommasInFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV with commas inside quoted fields
	csvContent := `NAME,ADDRESS,DESCRIPTION
"Smith, John","123 Main St, Apt 4B, City","Address with commas, lots of them"
"Doe, Jane","456 Oak Ave, Suite 200","Another address, with commas"`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":        "Smith, John",
			"ADDRESS":     "123 Main St, Apt 4B, City",
			"DESCRIPTION": "Address with commas, lots of them",
		},
		{
			"NAME":        "Doe, Jane",
			"ADDRESS":     "456 Oak Ave, Suite 200",
			"DESCRIPTION": "Another address, with commas",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataNewlinesInFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV with newlines inside quoted fields
	csvContent := `NAME,DESCRIPTION,NOTES
"John Doe","This is a long
description that spans
multiple lines","Single line note"
"Jane Smith","Another multi-line
description here","Another note"`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":        "John Doe",
			"DESCRIPTION": "This is a long\ndescription that spans\nmultiple lines",
			"NOTES":       "Single line note",
		},
		{
			"NAME":        "Jane Smith",
			"DESCRIPTION": "Another multi-line\ndescription here",
			"NOTES":       "Another note",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataMixedQuotingStyles(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV with mixed quoted and unquoted fields
	csvContent := `NAME,AGE,EMAIL,NOTES
"John Doe",30,john@example.com,"Has quotes in name"
Jane Smith,25,"jane@example.com","Normal entry"
"Bob Johnson",35,bob@test.com,No quotes in notes`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":  "John Doe",
			"AGE":   "30",
			"EMAIL": "john@example.com",
			"NOTES": "Has quotes in name",
		},
		{
			"NAME":  "Jane Smith",
			"AGE":   "25",
			"EMAIL": "jane@example.com",
			"NOTES": "Normal entry",
		},
		{
			"NAME":  "Bob Johnson",
			"AGE":   "35",
			"EMAIL": "bob@test.com",
			"NOTES": "No quotes in notes",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataEmptyFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV with empty fields
	csvContent := `NAME,EMAIL,PHONE,NOTES
"John Doe",john@example.com,,"No phone number"
"Jane Smith",,555-1234,"No email"
"Bob Johnson","","","All contact info missing"
"Alice Brown",alice@test.com,555-5678,`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":  "John Doe",
			"EMAIL": "john@example.com",
			"PHONE": "",
			"NOTES": "No phone number",
		},
		{
			"NAME":  "Jane Smith",
			"EMAIL": "",
			"PHONE": "555-1234",
			"NOTES": "No email",
		},
		{
			"NAME":  "Bob Johnson",
			"EMAIL": "",
			"PHONE": "",
			"NOTES": "All contact info missing",
		},
		{
			"NAME":  "Alice Brown",
			"EMAIL": "alice@test.com",
			"PHONE": "555-5678",
			"NOTES": "",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReadCSVDataSpecialJSONContent(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// CSV containing JSON-like content that might be used in curl templates
	csvContent := `NAME,JSON_DATA,API_KEY
"User1","{""name"": ""John"", ""age"": 30}","key-123-abc"
"User2","{""name"": ""Jane"", ""status"": ""active""}","key-456-def"`

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCSVData(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCSVData failed: %v", err)
	}

	expected := []map[string]string{
		{
			"NAME":      "User1",
			"JSON_DATA": `{"name": "John", "age": 30}`,
			"API_KEY":   "key-123-abc",
		},
		{
			"NAME":      "User2",
			"JSON_DATA": `{"name": "Jane", "status": "active"}`,
			"API_KEY":   "key-456-def",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
