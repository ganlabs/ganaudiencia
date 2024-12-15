package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/tebeka/selenium"
)

func GetWebdriver() (selenium.WebDriver, error) {
	var caps selenium.Capabilities
	var driver selenium.WebDriver

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

	switch environment {
	case "local":
		var driverport = GenerateRandomPort(999, 12000)
		opts := []selenium.ServiceOption{}
		driverPath, err := ExtractDriver()
		if err != nil {
			return nil, fmt.Errorf("erro ao extrair o driver: %w", err)
		}
		service, err := selenium.NewChromeDriverService(driverPath, driverport, opts...)
		if err != nil {
			return nil, fmt.Errorf("erro ao iniciar o serviço do ChromeDriver: %w", err)
		}
		log.Println(service)
		seleniumUrl = fmt.Sprintf("http://localhost:%d/wd/hub", driverport)
		driver, err = selenium.NewRemote(caps, seleniumUrl)
		if err != nil {
			return nil, fmt.Errorf("erro ao conectar ao WebDriver: %w", err)
		}
	case "docker":
		seleniumUrl = "http://selenium-hub:4444/wd/hub"
		driver, err = selenium.NewRemote(caps, seleniumUrl)
	default:
		seleniumUrl = "http://localhost:4444/wd/hub"
		driver, err = selenium.NewRemote(caps, seleniumUrl)
	}
	log.Println("WebDriver URL:", seleniumUrl)

	return driver, err
}
