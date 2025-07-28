package main

import (
	"reflect"
	"testing"
)

func TestSplitCurlCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []string
		hasError bool
	}{
		{
			name:     "Simple GET command",
			command:  `curl -X GET https://api.example.com`,
			expected: []string{"curl", "-X", "GET", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "POST with quoted header",
			command:  `curl -X POST -H "Content-Type: application/json" https://api.example.com`,
			expected: []string{"curl", "-X", "POST", "-H", "Content-Type: application/json", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "POST with single quoted data",
			command:  `curl -X POST -d '{"name": "test"}' https://api.example.com`,
			expected: []string{"curl", "-X", "POST", "-d", `{"name": "test"}`, "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Mixed quotes",
			command:  `curl -H "Authorization: Bearer token" -d '{"key": "value"}' https://api.example.com`,
			expected: []string{"curl", "-H", "Authorization: Bearer token", "-d", `{"key": "value"}`, "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Complex JSON data",
			command:  `curl -X POST -H "Content-Type: application/json" -d '{"name": "田中太郎", "email": "tanaka@example.com", "age": 30}' https://api.example.com`,
			expected: []string{"curl", "-X", "POST", "-H", "Content-Type: application/json", "-d", `{"name": "田中太郎", "email": "tanaka@example.com", "age": 30}`, "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Multiple headers",
			command:  `curl -H "Authorization: Bearer token" -H "Content-Type: application/json" https://api.example.com`,
			expected: []string{"curl", "-H", "Authorization: Bearer token", "-H", "Content-Type: application/json", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "URL with query parameters",
			command:  `curl "https://api.example.com/users?name=test&age=30"`,
			expected: []string{"curl", "https://api.example.com/users?name=test&age=30"},
			hasError: false,
		},
		{
			name:     "Escaped quotes in data",
			command:  `curl -d '{"message": "He said \"Hello\""}' https://api.example.com`,
			expected: []string{"curl", "-d", `{"message": "He said "Hello""}`, "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Unclosed double quote",
			command:  `curl -H "Content-Type: application/json https://api.example.com`,
			expected: nil,
			hasError: true,
		},
		{
			name:     "Unclosed single quote",
			command:  `curl -d '{"name": "test" https://api.example.com`,
			expected: nil,
			hasError: true,
		},
		{
			name:     "Empty command",
			command:  ``,
			expected: nil,
			hasError: false,
		},
		{
			name:     "Only curl",
			command:  `curl`,
			expected: []string{"curl"},
			hasError: false,
		},
		{
			name:     "Extra spaces",
			command:  `curl   -X   POST   https://api.example.com`,
			expected: []string{"curl", "-X", "POST", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Tabs and spaces",
			command:  "curl\t-X\tPOST\t\thttps://api.example.com",
			expected: []string{"curl", "-X", "POST", "https://api.example.com"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := splitCurlCommand(tt.command)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSplitCurlCommandSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []string
	}{
		{
			name:     "Unicode characters",
			command:  `curl -d '{"name": "田中太郎", "emoji": "🌸"}' https://api.example.com`,
			expected: []string{"curl", "-d", `{"name": "田中太郎", "emoji": "🌸"}`, "https://api.example.com"},
		},
		{
			name:     "Special symbols",
			command:  `curl -d '{"symbols": "@#$%^&*()"}' https://api.example.com`,
			expected: []string{"curl", "-d", `{"symbols": "@#$%^&*()"}`, "https://api.example.com"},
		},
		{
			name:     "Newline in quoted string",
			command:  "curl -d '{\"message\": \"line1\\nline2\"}' https://api.example.com",
			expected: []string{"curl", "-d", `{"message": "line1nline2"}`, "https://api.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := splitCurlCommand(tt.command)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSplitCurlCommandEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []string
		hasError bool
	}{
		{
			name:     "Quote at end",
			command:  `curl -H "Content-Type: application/json"`,
			expected: []string{"curl", "-H", "Content-Type: application/json"},
			hasError: false,
		},
		{
			name:     "Quote at start",
			command:  `"curl" -X GET https://api.example.com`,
			expected: []string{"curl", "-X", "GET", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Empty quotes",
			command:  `curl -H "" https://api.example.com`,
			expected: []string{"curl", "-H", "", "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Nested quotes different types",
			command:  `curl -d '{"message": "He said \"Hello\""}' https://api.example.com`,
			expected: []string{"curl", "-d", `{"message": "He said "Hello""}`, "https://api.example.com"},
			hasError: false,
		},
		{
			name:     "Backslash escape",
			command:  `curl -d '{"path": "C:\\Users\\test"}' https://api.example.com`,
			expected: []string{"curl", "-d", `{"path": "C:\Users\test"}`, "https://api.example.com"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := splitCurlCommand(tt.command)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
