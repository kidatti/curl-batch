package main

import (
	"os"
	"regexp"
	"strings"
)

// readCurlTemplate reads a curl template file and returns its content
func readCurlTemplate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// replaceTemplate replaces variables in the template with values from data
// Variables in the template should be in the format ${VARIABLE_NAME}
func (cb *CurlBatch) replaceTemplate(template string, data map[string]string) string {
	result := template
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	result = re.ReplaceAllStringFunc(result, func(match string) string {
		key := match[2 : len(match)-1] // Remove ${ and }
		if value, exists := data[key]; exists {
			return value
		}
		return match // Return unchanged if key doesn't exist
	})

	return result
}
