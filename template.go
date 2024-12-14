package main

type Processo struct {
	Processo      string   `json:"processo"`
	Classe        string   `json:"classe"`
	AudienciaData string   `json:"audiencia_data"`
	AudienciaHora string   `json:"audiencia_hora"`
	Valida        bool     `json:"valida"`
	Movimento     []string `json:"movimento"`
}
