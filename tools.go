package main

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func ValidateFormat(input string) (string, error) {
	// Remove all punctuation for validation without formatting
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "-", ""), ".", ""), " ", "")

	// Regular expression to match the unformatted format "XXXXXXXXXYYYY8ZZZZ"
	formatRegex := `^\d{7}\d{2}\d{4}8{1}19{1}\d{4}$`

	// Check if the cleaned input matches the format
	matched, err := regexp.MatchString(formatRegex, cleanInput)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", errors.New("invalid format")
	}

	// Check the second character is '8'
	if len(cleanInput) < 2 || cleanInput[1] != '8' {
		return "", errors.New("the second character must be '8'")
	}

	// Reformat the string with punctuation
	formattedInput := fmt.Sprintf("%s-%s.%s.8.19.%s", cleanInput[:7], cleanInput[7:9], cleanInput[9:13], cleanInput[16:])

	return formattedInput, nil
}

func GenerateRandomPort() int {
	// Create a new random generator with a seed based on the current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(3000) + 12000
}
