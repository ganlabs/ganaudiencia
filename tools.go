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
	// 1. Remover todas as pontuações: '-', '.', espaços
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "-", ""), ".", ""), " ", "")

	// 2. Verificar se a entrada limpa possui exatamente 20 dígitos
	if len(cleanInput) != 20 {
		return "", errors.New("entrada inválida: o tamanho deve ser exatamente 20 dígitos")
	}

	// 3. Verificar se todos os caracteres são dígitos
	if matched, _ := regexp.MatchString(`^\d{20}$`, cleanInput); !matched {
		return "", errors.New("entrada inválida: deve conter apenas dígitos")
	}

	// 4. Validar Segmento 3 (dígitos 9 a 12) >= 2000
	seg3 := cleanInput[9:13]
	seg3Int, err := strconv.Atoi(seg3)
	if err != nil {
		return "", errors.New("entrada inválida: o terceiro segmento deve ser numérico")
	}
	if seg3Int < 2000 {
		return "", errors.New("entrada inválida: o terceiro segmento deve ser >= 2000")
	}

	// 5. Validar Segmento 4 (dígito 13) seja '4' ou '8'
	seg4 := cleanInput[13:14]
	if seg4 != "4" && seg4 != "8" {
		return "", errors.New("entrada inválida: o quarto segmento deve ser '4' ou '8'")
	}

	// 6. Validar Segmento 5 (dígitos 14 e 15) conforme Segmento 4
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

	// 7. Segmento 6 (dígitos 16 a 19) já está validado como dígitos acima

	// 8. Reformatar a string com pontuação
	formattedInput := fmt.Sprintf(
		"%s-%s.%s.%s.%s.%s",
		cleanInput[:7],    // Segmento 1: primeiros 7 dígitos
		cleanInput[7:9],   // Segmento 2: dígitos 8 e 9
		cleanInput[9:13],  // Segmento 3: dígitos 10 a 13
		cleanInput[13:14], // Segmento 4: dígito 14
		cleanInput[14:16], // Segmento 5: dígitos 15 e 16
		cleanInput[16:20], // Segmento 6: dígitos 17 a 20
	)

	log.Println("Entrada formatada:", formattedInput)

	return formattedInput, nil
}

func GenerateRandomPort(quantity int, start int) int {
	// Create a new random generator with a seed based on the current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(quantity) + start
}

func ExtractDriver() (string, error) {
	// Determinar o nome do driver com base no sistema operacional
	var chromedriverPath string
	if runtime.GOOS == "windows" {
		chromedriverPath = "driver/chromedriver.exe"
	} else {
		chromedriverPath = "driver/chromedriver"
	}

	// Ler o conteúdo do binário embutido
	content, err := Chromedriver.ReadFile(chromedriverPath)
	if err != nil {
		return "", err
	}

	// Criar um diretório temporário para o driver
	tempDir, err := os.MkdirTemp("", "chromedriver-*")
	if err != nil {
		return "", err
	}

	// Caminho completo do arquivo no diretório temporário
	var tempFileName string
	if runtime.GOOS == "windows" {
		tempFileName = "chromedriver.exe"
	} else {
		tempFileName = "chromedriver"
	}
	tempFilePath := filepath.Join(tempDir, tempFileName)

	// Escrever o conteúdo no arquivo temporário
	if err := os.WriteFile(tempFilePath, content, 0755); err != nil {
		return "", err
	}

	return tempFilePath, nil
}
func ScraperDispatcher(lawsuit string) (Scraper, error) {
	lawsuit = strings.TrimSpace(lawsuit) // Remove espaços em branco
	switch {
	case Match(`^08[0-9]{5}-[0-9]{2}\.20[0-9]{2}\.8\.19\.[0-9]{4}$`, lawsuit):
		log.Println("PJE-RJ")
		return NewPjeRJ(), nil
	case Match(`^0[0-9]{6}-[0-9]{2}\.20[0-9]{2}\.8\.26\.[0-9]{4}$`, lawsuit):
		log.Println("ESAJ-SP")
		return NewEsajSP(), nil
	case Match(`^0[01][0-9]{7}20[0-9]{2}819[0-9]{4}$`, lawsuit):
		log.Println("TJRJ")
		return nil, nil
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
