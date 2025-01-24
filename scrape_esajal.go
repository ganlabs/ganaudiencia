package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type EsajAL struct {
}

func NewEsajAL() Scraper {
	return &EsajAL{}
}

func (s *EsajAL) Scrape(lawsuit string) (Hearing, error) {

	url := "https://www2.tjal.jus.br/cpopg/open.do"

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

	if err == nil {
		expand.SendKeys("\n")
	}

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
func (s *EsajAL) ValidateDate(date string) bool {
	d, err := time.Parse("02/01/2006", date)

	if err != nil {
		return false
	}

	if d.After(time.Now()) {
		return true
	}

	return false
}
