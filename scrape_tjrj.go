package main

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

type Tjrj struct {
}

func NewTjrj() Scraper {
	return &Tjrj{}
}

func (s *Tjrj) Scrape(lawsuit string) (Hearing, error) {

	url := fmt.Sprintf("https://www3.tjrj.jus.br/consultaprocessual/#/consultapublica?numProcessoCNJ=%s", lawsuit)

	hr := Hearing{
		Lawsuit: lawsuit,
	}

	driver, err := GetWebdriver()
	if err != nil {
		return hr, fmt.Errorf("erro ao obter o driver: %w", err)
	}

	defer driver.Quit()

	err = driver.Get(url)
	if err != nil {
		return hr, fmt.Errorf("erro ao acessar a página: %w", err)
	}

	err = waitForPageLoad(driver, 5*time.Second)
	if err != nil {
		return hr, fmt.Errorf("erro ao aguardar a página carregar: %w", err)
	}

	frame, err := driver.FindElements(selenium.ByTagName, "iframe")
	if err != nil {
		return hr, fmt.Errorf("erro ao localizar o iframe: %w", err)
	}

	fmt.Println(frame)

	movementsButton, err := driver.FindElement(selenium.ByXPATH, "/html/body/app-root/app-detalhes-processo/section/div/div/div[1]/div[2]/button[2]")
	if err != nil {
		return hr, fmt.Errorf("erro ao localizar o botão de movimentos: %w", err)
	}

	_, err = movementsButton.IsDisplayed()
	if err != nil {
		return hr, fmt.Errorf("erro ao verificar se o botão de movimentos foi exibido: %w", err)
	}

	fmt.Println(movementsButton.Text())

	time.Sleep(5 * time.Second)

	return Hearing{}, nil
}

func waitForPageLoad(driver selenium.WebDriver, timeout time.Duration) error {
	start := time.Now()
	for {
		state, err := driver.ExecuteScript("return document.readyState", nil)
		if err != nil {
			return fmt.Errorf("erro ao executar script: %v", err)
		}

		if state == "complete" {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout ao carregar a página")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
