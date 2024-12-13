package main

type Hearing struct {
	Lawsuit     string `json:"processo"`
	Class       string `json:"classe"`
	HearingDate string `json:"audiencia_data"`
	HearingTime string `json:"audiencia_hora"`
	IsValid     bool   `json:"valida"`
	Movement    string `json:"movimento"`
}
