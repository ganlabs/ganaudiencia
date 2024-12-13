package main

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateFormat validates the input format and conditions
func ValidateFormat(input string) (string, error) {
	// Remove all punctuation for validation without formatting
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "-", ""), ".", ""), " ", "")

	// Regular expression to match the unformatted format "XXXXXXXXXYYYY8ZZZZ"
	formatRegex := `^\d{7}\d{2}\d{4}8{1}19{1}\d{4}$`

	// Check if the cleaned input matches the format
	matched, err := regexp.MatchString(formatRegex, cleanInput)
	if err != nil {
		return cleanInput, err
	}
	if !matched {
		return cleanInput, errors.New("invalid format")
	}

	// Check the second character is '8'
	if len(cleanInput) < 2 || cleanInput[1] != '8' {
		return cleanInput, errors.New("the second character must be '8'")
	}

	return cleanInput, nil
}
