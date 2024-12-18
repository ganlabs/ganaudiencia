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

	// movement, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/table[7]/tbody/tr/td")
	movement, err := driver.FindElement(selenium.ByXPATH, "//*[@id=\"processoSemAudiencias\"]")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o movimento: %w", err)
	}

	movementText, err := movement.Text()
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao obter o texto do movimento: %w", err)
	}

	var hearingDate string = "01/01/1900"
	var hearingTime string = "00:00"

	if strings.Contains(strings.ToUpper(movementText), "AUDIÊNCIA") && strings.Contains(strings.ToUpper(movementText), "DESIGNADA") {
		hearingDate = detectAndFormatFirstDate(movementText)
		hearingTime = "00:00"
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
