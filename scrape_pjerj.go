package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type PjeRJ struct {
}

func NewPjeRJ() Scraper {
	return &PjeRJ{}
}

func (s *PjeRJ) Scrape(lawsuit string) (Hearing, error) {
	mv, err := s.FetchMovements(lawsuit)
	if err != nil {
		log.Println(err)
		return Hearing{}, err
	}

	hd, ht := s.ExtractHearingDates(mv)
	log.Println("Hearing date/time:", hd, ht)

	hr := Hearing{
		Lawsuit:     lawsuit,
		Class:       s.ExtractClass(mv[0]),
		HearingDate: hd,
		HearingTime: ht,
		IsValid:     s.ValidateDate(hd),
		Movement:    mv,
	}

	return hr, nil
}

func (s *PjeRJ) FetchMovements(processNumber string) ([]string, error) {

	driver, err := GetWebdriver()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter o driver: %w", err)
	}

	defer driver.Quit()
	defer driver.Close()

	var processLink selenium.WebElement
	attempts := 0
	for attempts < 10 {
		err = driver.Get("https://tjrj.pje.jus.br/1g/ConsultaPublica/listView.seam")
		if err != nil {
			return nil, fmt.Errorf("erro ao acessar a página: %w", err)
		}

		input, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/form/div[1]/div/div/div/div/div[1]/div/div[2]/input")
		if err != nil {
			return nil, fmt.Errorf("erro ao localizar o campo de entrada: %w", err)
		}
		if err := input.SendKeys(processNumber); err != nil {
			return nil, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
		}

		time.Sleep(1 * time.Second)

		searchButton, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/form/div[1]/div/div/div/div/div[8]/div/input")
		if err != nil {
			return nil, fmt.Errorf("erro ao localizar o botão de pesquisa: %w", err)
		}
		if err := searchButton.Click(); err != nil {
			return nil, fmt.Errorf("erro ao clicar no botão de pesquisa: %w", err)
		}

		time.Sleep(1 * time.Second)

		processLink, err = driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/form/div[2]/div/table/tbody/tr/td[2]/a/b")
		if err == nil {
			if err := processLink.Click(); err != nil {
				return nil, fmt.Errorf("erro ao clicar no link do processo: %w", err)
			}
			break
		}

		attempts++
		if attempts == 20 {
			driver.Quit()
			fmt.Printf("falha ao encontrar o processo após 20 tentativas: %s", processNumber)
			return []string{}, nil
		}
	}

	windows, err := driver.WindowHandles()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter as janelas: %w", err)
	}
	if err := driver.SwitchWindow(windows[1]); err != nil {
		return nil, fmt.Errorf("erro ao trocar para a nova janela: %w", err)
	}

	time.Sleep(1 * time.Second)

	currentURL, err := driver.CurrentURL()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter a URL atual: %w", err)
	}
	if err := driver.Get(currentURL); err != nil {
		return nil, fmt.Errorf("erro ao acessar a URL do processo: %w", err)
	}

	time.Sleep(1 * time.Second)

	var movements []string

	class, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/table/tbody/tr[2]/td/table/tbody/tr/td/form/div/div[1]/div[3]/table/tbody/tr[1]/td[3]/span/div/div[2]")
	if err != nil {
		return nil, fmt.Errorf("erro ao obter o texto da classe: %w", err)
	}

	classText, _ := class.Text()

	movements = append(movements, classText)

	err = s.FetchMovementsFromPage(driver, &movements)
	if err != nil {
		return nil, err
	}

	page := 2
	for {
		pageInput, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/table/tbody/tr[2]/td/table/tbody/tr/td/div[5]/div[2]/div/form/table/tbody/tr[1]/td[5]/input")
		if err != nil {
			break
		}

		if err := pageInput.Clear(); err != nil {
			return nil, fmt.Errorf("erro ao limpar o campo de entrada: %w", err)
		}

		if err := pageInput.SendKeys(fmt.Sprintf("%d", page)); err != nil {
			return nil, fmt.Errorf("erro ao enviar texto para o campo de entrada: %w", err)
		}

		if err := pageInput.SendKeys(selenium.EnterKey); err != nil {
			return nil, fmt.Errorf("erro ao pressionar Enter: %w", err)
		}

		time.Sleep(1 * time.Second)

		_, err = driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/table/tbody/tr[2]/td/table/tbody/tr/td/div[5]/div[2]/table/tbody")
		if err != nil {
			break
		}

		err = s.FetchMovementsFromPage(driver, &movements)
		if err != nil {
			return nil, err
		}

		page++

		nextPageIndicator, err := driver.FindElement(selenium.ByXPATH, "//span[contains(@class, 'rich-datascr-act') and text()='"+fmt.Sprintf("%d", page-1)+"']")
		if err != nil || nextPageIndicator == nil {
			break
		}
	}
	return movements, nil
}

func (s *PjeRJ) FetchMovementsFromPage(driver selenium.WebDriver, movements *[]string) error {
	table, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div/div/div/div[2]/table/tbody/tr[2]/td/table/tbody/tr/td/div[5]/div[2]/table/tbody")
	if err != nil {
		return fmt.Errorf("erro ao localizar a tabela de movimentos: %w", err)
	}
	rows, err := table.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return fmt.Errorf("erro ao localizar as linhas da tabela: %w", err)
	}
	for _, row := range rows {
		text, err := row.Text()
		if err != nil {
			log.Printf("erro ao obter o texto da linha: %v", err)
			continue
		}
		if text != "" {
			*movements = append(*movements, text)
		}
	}
	return nil
}

func (s *PjeRJ) ExtractHearingDates(lines []string) (date string, time string) {
	for _, line := range lines {
		line = strings.ToUpper(line)

		if parts := strings.SplitN(line, " - ", 2); len(parts) > 1 {
			line = parts[1]
		}

		if strings.Contains(line, "AUDIÊNCIA") && strings.Contains(line, "REALIZADA") {
			break
		}
		if strings.Contains(line, "AUDIÊNCIA") && strings.Contains(line, "CANCELADA") {
			break
		}
		// if strings.Contains(line, "AUDIÊNCIA") && strings.Contains(line, "REDESIGNADA") {
		// 	break
		// }
		if strings.Contains(line, "AUDIÊNCIA") && strings.Contains(line, "DESIGNADA") {
			regexDate := regexp.MustCompile(`(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[012])/(19|20)\d{2}`)
			regexTime := regexp.MustCompile(`\b(\d{2}:\d{2})\b`)

			dateMatch := regexDate.FindString(line)
			timeMatch := regexTime.FindString(line)

			if dateMatch != "" && timeMatch != "" {
				return dateMatch, timeMatch
			}

			break
		}
	}

	return "01/01/1900", "00:00"
}

func (s *PjeRJ) ExtractClass(text string) string {
	text = strings.ToUpper(text)
	if strings.Contains(text, "JUIZADO ESPECIAL") {
		return "JEC"
	} else if strings.Contains(text, "COMUM CÍVEL") {
		return "VC"
	}
	return text
}

func (s *PjeRJ) ValidateDate(date string) bool {
	d, err := time.Parse("02/01/2006", date)

	if err != nil {
		return false
	}

	if d.After(time.Now()) {
		return true
	}

	return false
}
