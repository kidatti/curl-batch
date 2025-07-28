package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// splitCurlCommand parses a curl command string and splits it into arguments
// while properly handling quoted strings (both single and double quotes)
func splitCurlCommand(command string) ([]string, error) {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	escaped := false

	for i := 0; i < len(command); i++ {
		char := command[i]

		if escaped {
			current.WriteByte(char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if !inQuotes && (char == '\'' || char == '"') {
			inQuotes = true
			quoteChar = char
			continue
		}

		if inQuotes && char == quoteChar {
			inQuotes = false
			quoteChar = 0
			// Always add the current content when closing quotes (even if empty)
			parts = append(parts, current.String())
			current.Reset()
			continue
		}

		if !inQuotes && (char == ' ' || char == '\t') {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(char)
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in command")
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts, nil
}

// executeRequest parses and executes a curl command
func (cb *CurlBatch) executeRequest(curlCommand string) (string, error) {
	parts, err := splitCurlCommand(curlCommand)
	if err != nil {
		return "", fmt.Errorf("failed to parse curl command: %w", err)
	}

	if len(parts) < 2 || parts[0] != "curl" {
		return "", fmt.Errorf("invalid curl command: %s", curlCommand)
	}

	var method, url, body string
	var headers []string

	for i := 1; i < len(parts); i++ {
		switch parts[i] {
		case "-X":
			if i+1 < len(parts) {
				method = parts[i+1]
				i++
			}
		case "-H":
			if i+1 < len(parts) {
				headers = append(headers, parts[i+1])
				i++
			}
		case "-d":
			if i+1 < len(parts) {
				body = parts[i+1]
				i++
			}
		default:
			if strings.HasPrefix(parts[i], "http") {
				url = parts[i]
			}
		}
	}

	if method == "" {
		method = "GET"
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return fmt.Sprintf("Status: %s\nHeaders: %v\nBody: %s", resp.Status, resp.Header, string(respBody)), nil
}
