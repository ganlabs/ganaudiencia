package main

type Scraper interface {
	Scrape(lawsuit string) (Hearing, error)
}

type Hearing struct {
	Lawsuit      string   `json:"processo"`
	Class        string   `json:"classe"`
	HearingDate  string   `json:"audiencia_data"`
	HearingTime  string   `json:"audiencia_hora"`
	IsValid      bool     `json:"valida"`
	Author       string   `json:"autor"`
	Jurisdiction string   `json:"vara"`
	Location     string   `json:"comarca"`
	Customer     string   `json:"reu"`
	Movement     []string `json:"movimento"`
}
