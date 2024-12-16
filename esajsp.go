package main

import (
	"fmt"

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

	input, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/section/div[2]/div/div[1]/div[1]/span[1]/input[1]")

	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o campo de entrada: %w", err)
	}
	if err := input.SendKeys(lawsuit); err != nil {
		return Hearing{}, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
	}

	searchButton, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/section/div[4]/div/input")
	if err != nil {
		return Hearing{}, fmt.Errorf("erro ao localizar o botão de pesquisa: %w", err)
	}
	if err := searchButton.Click(); err != nil {
		return Hearing{}, fmt.Errorf("erro ao clicar no botão de pesquisa: %w", err)
	}

	return Hearing{}, nil

}
