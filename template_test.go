package main

import (
	"os"
	"strings"
	"testing"
)

func TestReadCurlTemplate(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "curl_template_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `curl -X POST -H "Content-Type: application/json" -d '{"name": "${NAME}"}' https://api.example.com`
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCurlTemplate(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCurlTemplate failed: %v", err)
	}

	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

func TestReadCurlTemplateWithWhitespace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "curl_template_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `  curl -X POST https://api.example.com  
`
	expected := `curl -X POST https://api.example.com`

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := readCurlTemplate(tmpFile.Name())
	if err != nil {
		t.Fatalf("readCurlTemplate failed: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestReadCurlTemplateFileNotExists(t *testing.T) {
	_, err := readCurlTemplate("nonexistent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, but got none")
	}
}

func TestReplaceTemplate(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -H "Content-Type: application/json" -d '{"name": "${NAME}", "email": "${EMAIL}", "age": ${AGE}}' https://api.example.com`
	data := map[string]string{
		"NAME":  "田中太郎",
		"EMAIL": "tanaka@example.com",
		"AGE":   "30",
	}

	result := cb.replaceTemplate(template, data)
	expected := `curl -X POST -H "Content-Type: application/json" -d '{"name": "田中太郎", "email": "tanaka@example.com", "age": 30}' https://api.example.com`

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateWithMissingVariables(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -d '{"name": "${NAME}", "missing": "${MISSING_VAR}"}' https://api.example.com`
	data := map[string]string{
		"NAME": "田中太郎",
	}

	result := cb.replaceTemplate(template, data)

	// Missing variables should remain unchanged
	if !strings.Contains(result, "${MISSING_VAR}") {
		t.Errorf("Expected missing variable to remain as ${MISSING_VAR}, got %q", result)
	}

	if !strings.Contains(result, "田中太郎") {
		t.Errorf("Expected existing variable to be replaced, got %q", result)
	}
}

func TestReplaceTemplateNoVariables(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X GET https://api.example.com/status`
	data := map[string]string{
		"NAME": "田中太郎",
	}

	result := cb.replaceTemplate(template, data)

	if result != template {
		t.Errorf("Expected template to remain unchanged when no variables, got %q", result)
	}
}

func TestReplaceTemplateEmptyData(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -d '{"name": "${NAME}"}' https://api.example.com`
	data := map[string]string{}

	result := cb.replaceTemplate(template, data)

	// Variables should remain unchanged when no data provided
	if !strings.Contains(result, "${NAME}") {
		t.Errorf("Expected variable to remain as ${NAME}, got %q", result)
	}
}

func TestReplaceTemplateMultipleSameVariable(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -d '{"name": "${NAME}", "fullname": "${NAME}"}' https://api.example.com/${NAME}`
	data := map[string]string{
		"NAME": "田中太郎",
	}

	result := cb.replaceTemplate(template, data)
	expected := `curl -X POST -d '{"name": "田中太郎", "fullname": "田中太郎"}' https://api.example.com/田中太郎`

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateSpecialCharacters(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -d '{"message": "${MESSAGE}"}' https://api.example.com`
	data := map[string]string{
		"MESSAGE": `Hello "World" & <script>alert('test')</script>`,
	}

	result := cb.replaceTemplate(template, data)
	expected := `curl -X POST -d '{"message": "Hello "World" & <script>alert('test')</script>"}' https://api.example.com`

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateUnicodeCharacters(t *testing.T) {
	cb := &CurlBatch{}

	template := `curl -X POST -d '{"name": "${NAME}", "emoji": "${EMOJI}"}' https://api.example.com`
	data := map[string]string{
		"NAME":  "田中太郎",
		"EMOJI": "🌸🗾🎌",
	}

	result := cb.replaceTemplate(template, data)
	expected := `curl -X POST -d '{"name": "田中太郎", "emoji": "🌸🗾🎌"}' https://api.example.com`

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateEdgeCases(t *testing.T) {
	cb := &CurlBatch{}

	tests := []struct {
		name     string
		template string
		data     map[string]string
		expected string
	}{
		{
			name:     "Empty template",
			template: "",
			data:     map[string]string{"NAME": "test"},
			expected: "",
		},
		{
			name:     "Variable at start",
			template: "${NAME} is here",
			data:     map[string]string{"NAME": "田中太郎"},
			expected: "田中太郎 is here",
		},
		{
			name:     "Variable at end",
			template: "Hello ${NAME}",
			data:     map[string]string{"NAME": "田中太郎"},
			expected: "Hello 田中太郎",
		},
		{
			name:     "Malformed variable",
			template: "Hello ${NAME",
			data:     map[string]string{"NAME": "田中太郎"},
			expected: "Hello ${NAME",
		},
		{
			name:     "Empty variable name",
			template: "Hello ${}",
			data:     map[string]string{"": "test"},
			expected: "Hello ${}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cb.replaceTemplate(tt.template, tt.data)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
