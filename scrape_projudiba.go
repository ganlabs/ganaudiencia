package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type ProjudiBA struct {
}

func NewProjudiBA() Scraper {
	return &ProjudiBA{}
}

func (s *ProjudiBA) Scrape(lawsuit string) (Hearing, error) {

	url := "https://projudi.tjba.jus.br/projudi/PaginaPrincipal.jsp"

	driver, err := GetWebdriver()
	defer driver.Quit()

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao obter o driver: %w", err)
	}

	err = driver.Get(url)
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao acessar a página: %w", err)
	}

	input1, err := driver.FindElement(selenium.ByID, "numeroProcesso")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o campo de entrada: %w", err)
	}

	if err := input1.SendKeys(lawsuit + "\n"); err != nil {
		return Hearing{}, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
	}

	time.Sleep(1 * time.Second)

	class, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[4]/table/tbody/tr[10]/td[2]")

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

	movement, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[6]/table/tbody")

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
		if strings.Contains(strings.ToUpper(line), "AGENDADA PARA") {
			hearingDate = detectAndFormatProjudiDate(line)
			fmt.Println("###########################################")
			fmt.Println(hearingDate)
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

func (s *ProjudiBA) ValidateDate(date string) bool {
	d, err := time.Parse("02/01/2006", date)

	if err != nil {
		return false
	}

	if d.After(time.Now()) {
		return true
	}

	return false
}

func detectAndFormatProjudiDate(text string) string {
	pattern := `\b(\d{2})\s+de\s+([a-zçã]+)\s+de\s+(\d{2,4})\b`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(strings.ToLower(text))

	if len(match) == 4 {
		day := match[1]
		monthText := match[2]
		year := match[3]

		month, exists := monthNames[monthText]
		if !exists {
			return "01/01/1900"
		}

		if len(year) == 2 {
			yearInt, err := strconv.Atoi(year)
			if err != nil {
				return "01/01/1900"
			}
			if yearInt <= 50 {
				year = fmt.Sprintf("20%02d", yearInt)
			} else {
				year = fmt.Sprintf("19%02d", yearInt)
			}
		}

		dateStr := fmt.Sprintf("%s/%s/%s", day, month, year)

		if isValidDate(dateStr) {
			return dateStr
		}
	}

	return "01/01/1900"
}
