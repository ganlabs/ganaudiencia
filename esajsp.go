package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type EsajSP struct {
}

func NewEsajSP() Scraper {
	return &EsajSP{}
}

func (s *EsajSP) Scrape(lawsuit string) (Hearing, error) {

	url := "https://esaj.tjsp.jus.br/cpopg/open.do"

	driver, err := GetWebdriver()
	defer driver.Quit()

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao obter o driver: %w", err)
	}

	err = driver.Get(url)
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao acessar a página: %w", err)
	}

	input1, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/section/div[2]/div/div[1]/div[1]/span[1]/input[1]")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o campo de entrada: %w", err)
	}

	input2, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/section/div[2]/div/div[1]/div[1]/span[1]/input[3]")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o campo de entrada: %w", err)
	}

	if err := input1.SendKeys(lawsuit[:15]); err != nil {
		return Hearing{}, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
	}

	if err := input2.SendKeys(lawsuit[21:]); err != nil {
		return Hearing{}, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
	}

	searchButton, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/section/div[4]/div/input")
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o botão de pesquisa: %w", err)
	}
	if err := searchButton.Click(); err != nil {
		return Hearing{}, fmt.Errorf("erro ao clicar no botão de pesquisa: %w", err)
	}

	time.Sleep(1 * time.Second)

	class, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[1]/div[2]/div/div[2]/div[1]/div/span")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar a classe: %w", err)
	}

	classText, err := class.Text()
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao obter o texto da classe: %w", err)
	}

	if strings.Contains(strings.ToUpper(classText), "JUIZADO ESPECIAL") {
		classText = "JEC"
	} else if strings.Contains(strings.ToUpper(classText), "COMUM CÍVEL") {
		classText = "VC"
	}

	expand, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/div[4]/a")

	//if err != nil {
	//	return Hearing{}, fmt.Errorf("erro ao selecionar movimentação: %w", err)
	//}

	if err == nil {
		err = expand.SendKeys("\n")
	}

	//if err != nil {
	//	return Hearing{}, fmt.Errorf("erro ao expandir movimentação: %w", err)
	//}

	//movement, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/table[7]/tbody/tr/td")
	movement, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/table[2]")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o movimento: %w", err)
	}

	movementText, err := movement.Text()
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao obter o texto do movimento: %w", err)
	}

	var hearingDate string = "01/01/1900"
	var hearingTime string = "00:00"

	lines := strings.Split(movementText, "\n")

	for _, line := range lines {
		if strings.Contains(strings.ToUpper(line), "DATA") && strings.Contains(strings.ToUpper(line), "HORA") && strings.Contains(strings.ToUpper(line), "SITUACÃO") {
			fmt.Println(line)
			hearingDate = detectAndFormatFirstDate(line)
			hearingTime = detectAndFormatFirstTime(line)
			break
		}
	}

	return Hearing{
		Lawsuit:     lawsuit,
		Class:       classText,
		HearingDate: hearingDate,
		HearingTime: hearingTime,
		IsValid:     s.ValidateDate(hearingDate),
		Movement:    []string{movementText},
	}, nil

}

var monthNames = map[string]string{
	"janeiro": "01", "fevereiro": "02", "março": "03", "abril": "04",
	"maio": "05", "junho": "06", "julho": "07", "agosto": "08",
	"setembro": "09", "outubro": "10", "novembro": "11", "dezembro": "12",
}

func detectAndFormatFirstDate(text string) string {
	// List of date regex patterns
	patterns := []string{
		`\b(\d{2})/(\d{2})/(\d{2,4})\b`,        // dd/mm/yy or dd/mm/yyyy
		`\b(\d{2}) de ([a-z]+) de (\d{2,4})\b`, // dd de mmmm de yy or yyyy
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(text)

		if len(match) == 4 {
			day := match[1]
			month := match[2]
			year := match[3]

			// Handle month as numeric or text
			if val, exists := monthNames[strings.ToLower(month)]; exists {
				month = val
			}

			// Handle 2-digit year to 4-digit year conversion
			if len(year) == 2 {
				yearInt, _ := strconv.Atoi(year)
				if yearInt <= 50 {
					year = fmt.Sprintf("20%02d", yearInt)
				} else {
					year = fmt.Sprintf("19%02d", yearInt)
				}
			}

			// Validate and return the date
			dateStr := fmt.Sprintf("%s/%s/%s", day, month, year)
			if isValidDate(dateStr) {
				return dateStr
			}
		}
	}

	return "01/01/1900"
}

// isValidDate validates if the date is a real date
func isValidDate(date string) bool {
	_, err := time.Parse("02/01/2006", date)
	return err == nil
}

// detectAndFormatFirstTime detecta a primeira ocorrência de hora no texto e a formata como "HH:MM"
// Se nenhuma hora válida for encontrada, retorna "00:00"
func detectAndFormatFirstTime(text string) string {
	// Lista de padrões regex para diferentes formatos de hora
	patterns := []string{
		`\b([01]?\d|2[0-3]):([0-5]\d)\b`,             // Formato 24 horas: HH:MM
		`\b([1-9]|1[0-2])\s?(AM|PM|am|pm)\b`,         // Formato 12 horas com AM/PM
		`\b([01]?\d|2[0-3])\s?horas?\s?([0-5]\d)?\b`, // Formato com "horas" opcionalmente seguido por minutos
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)

		if len(matches) >= 3 {
			var hour, minute string

			// Formato 24 horas: HH:MM
			if len(matches) == 3 && strings.Contains(pattern, ":") {
				hour = matches[1]
				minute = matches[2]
			} else if len(matches) >= 2 && strings.ContainsAny(pattern, "AMPMampm") {
				// Formato 12 horas com AM/PM
				hour = matches[1]
				minute = "00" // Assume minutos como "00" se não especificados

				// Opcional: Converter para 24 horas se AM/PM estiver presente
				if len(matches) == 3 {
					meridiem := strings.ToLower(matches[2])
					hourInt, _ := strconv.Atoi(hour)
					if meridiem == "pm" && hourInt != 12 {
						hourInt += 12
					} else if meridiem == "am" && hourInt == 12 {
						hourInt = 0
					}
					hour = fmt.Sprintf("%02d", hourInt)
				}
			} else if len(matches) >= 2 {
				// Formato com "horas" e possivelmente minutos
				hour = matches[1]
				if len(matches) == 3 && matches[2] != "" {
					minute = matches[2]
				} else {
					minute = "00"
				}
			}

			// Adicionar zeros à esquerda se necessário
			hour = fmt.Sprintf("%02s", hour)
			minute = fmt.Sprintf("%02s", minute)

			// Validar e retornar a hora formatada
			timeStr := fmt.Sprintf("%s:%s", hour, minute)
			if isValidTime(timeStr) {
				return timeStr
			}
		}
	}

	return "00:00"
}

// isValidTime verifica se a string fornecida está no formato de hora válido "HH:MM"
func isValidTime(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}

	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return false
	}

	if hour < 0 || hour > 23 {
		return false
	}

	if minute < 0 || minute > 59 {
		return false
	}

	return true
}

func (s *EsajSP) ValidateDate(date string) bool {
	d, err := time.Parse("02/01/2006", date)

	if err != nil {
		return false
	}

	// compare the date with the current date
	if d.After(time.Now()) {
		return true
	}

	return false
}
