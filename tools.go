package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ValidateFormat(input string) (string, error) {

	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "-", ""), ".", ""), " ", "")

	if len(cleanInput) != 20 {
		return "", errors.New("entrada inválida: o tamanho deve ser exatamente 20 dígitos")
	}

	if matched, _ := regexp.MatchString(`^\d{20}$`, cleanInput); !matched {
		return "", errors.New("entrada inválida: deve conter apenas dígitos")
	}

	seg3 := cleanInput[9:13]
	seg3Int, err := strconv.Atoi(seg3)
	if err != nil {
		return "", errors.New("entrada inválida: o terceiro segmento deve ser numérico")
	}
	if seg3Int < 2000 {
		return "", errors.New("entrada inválida: o terceiro segmento deve ser >= 2000")
	}

	seg4 := cleanInput[13:14]
	if seg4 != "4" && seg4 != "8" {
		return "", errors.New("entrada inválida: o quarto segmento deve ser '4' ou '8'")
	}

	seg5 := cleanInput[14:16]
	seg5Int, err := strconv.Atoi(seg5)
	if err != nil {
		return "", errors.New("entrada inválida: o quinto segmento deve ser numérico")
	}

	if seg4 == "4" {
		if seg5Int < 1 || seg5Int > 5 {
			return "", errors.New("entrada inválida: o quinto segmento deve estar entre 01 e 05 quando o quarto segmento for '4'")
		}
	} else if seg4 == "8" {
		if seg5Int < 1 || seg5Int > 27 {
			return "", errors.New("entrada inválida: o quinto segmento deve estar entre 01 e 27 quando o quarto segmento for '8'")
		}
	}

	formattedInput := fmt.Sprintf(
		"%s-%s.%s.%s.%s.%s",
		cleanInput[:7],
		cleanInput[7:9],
		cleanInput[9:13],
		cleanInput[13:14],
		cleanInput[14:16],
		cleanInput[16:20],
	)

	log.Println("Entrada formatada:", formattedInput)

	return formattedInput, nil
}

func GenerateRandomPort(quantity int, start int) int {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(quantity) + start
}

func ExtractDriver() (string, error) {

	var chromedriverPath string
	if runtime.GOOS == "windows" {
		chromedriverPath = "driver/chromedriver.exe"
	} else {
		chromedriverPath = "driver/chromedriver"
	}

	content, err := Chromedriver.ReadFile(chromedriverPath)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "chromedriver-*")
	if err != nil {
		return "", err
	}

	var tempFileName string
	if runtime.GOOS == "windows" {
		tempFileName = "chromedriver.exe"
	} else {
		tempFileName = "chromedriver"
	}
	tempFilePath := filepath.Join(tempDir, tempFileName)

	if err := os.WriteFile(tempFilePath, content, 0755); err != nil {
		return "", err
	}

	return tempFilePath, nil
}
func ScraperDispatcher(lawsuit string) (Scraper, error) {
	lawsuit = strings.TrimSpace(lawsuit)
	switch {
	case Match(`^0[8-9]{1}[0-9]{5}-[0-9]{2}\.20[0-9]{2}\.8\.19\.[0-9]{4}$`, lawsuit):
		log.Println("PJE-RJ")
		return NewPjeRJ(), nil
	case Match(`^[0-9]{7}-[0-9]{2}\.20[0-9]{2}\.8\.26\.[0-9]{4}$`, lawsuit):
		log.Println("ESAJ-SP")
		return NewEsajSP(), nil
	case Match(`^0[0-1]{1}[0-9]{5}-[0-9]{2}\.20[0-9]{2}\.8\.19\.[0-9]{4}$`, lawsuit):
		log.Println("TJRJ")
		return NewTjrj(), nil
	default:
		return nil, errors.New("invalid lawsuit")
	}
}

func Match(pattern, input string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Erro ao compilar regex: %v\n", err)
		return false
	}
	match := re.MatchString(input)
	fmt.Printf("Matching input '%s' against pattern '%s': %v\n", input, pattern, match)
	return match
}
