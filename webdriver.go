package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/tebeka/selenium"
)

func GetWebdriver() (selenium.WebDriver, error) {
	var caps selenium.Capabilities
	var webdriver selenium.WebDriver

	switch runtime.GOOS {
	case "windows":
		caps = selenium.Capabilities{
			"browserName": "chrome",
			"goog:chromeOptions": map[string]interface{}{
				"args": []string{
					"--no-sandbox",
					"--disable-dev-shm-usage",
				},
			},
		}
	default:
		caps = selenium.Capabilities{
			"browserName": "chrome",
			"goog:chromeOptions": map[string]interface{}{
				"args": []string{
					"--no-sandbox",
					"--disable-dev-shm-usage",
				},
			},
		}
	}

	var err error
	var seleniumUrl string

	switch Environment {
	case "docker":
		seleniumUrl = "http://selenium-hub:4444/wd/hub"
		webdriver, err = selenium.NewRemote(caps, seleniumUrl)
	default:
		var driverport = GenerateRandomPort(999, 12000)
		driverPath, err := ExtractDriver()

		if err != nil {
			return nil, fmt.Errorf("erro ao extrair o driver: %w", err)
		}

		service, err := selenium.NewChromeDriverService(driverPath, driverport)
		if err != nil {
			return nil, fmt.Errorf("erro ao iniciar o serviço do ChromeDriver: %w", err)
		}
		log.Printf("ChromeDriver iniciado na porta %d", driverport)
		log.Printf("Acessar o navegador em http://localhost:%d/wd/hub", driverport)
		log.Println("Para sair, pressione CTRL+C")
		log.Println(service)

		seleniumUrl = fmt.Sprintf("http://localhost:%d/wd/hub", driverport)
		webdriver, err = selenium.NewRemote(caps, seleniumUrl)
		if err != nil {
			return nil, fmt.Errorf("erro ao conectar ao WebDriver: %w", err)
		}
	}

	return webdriver, err
}
