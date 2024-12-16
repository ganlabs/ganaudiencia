package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Casos Válidos
		{
			name:        "Formato válido com pontuação - seg4=4, seg5=03",
			input:       "1234567-89.2001.4.03.5678",
			expected:    "1234567-89.2001.4.03.5678",
			expectError: false,
		},
		{
			name:        "Formato válido com pontuação - seg4=8, seg5=15",
			input:       "7654321-98.2020.8.15.8765",
			expected:    "7654321-98.2020.8.15.8765",
			expectError: false,
		},
		{
			name:        "Formato válido sem pontuação - seg4=4, seg5=05",
			input:       "11111112220004 05 4444",
			expected:    "1111111-22.2000.4.05.4444",
			expectError: false,
		},
		{
			name:        "Formato válido sem pontuação - seg4=8, seg5=27",
			input:       "66666667720308 27 9999",
			expected:    "6666666-77.2030.8.27.9999",
			expectError: false,
		},
		{
			name:        "Caso de borda - seg3=2000",
			input:       "2222222-33.2000.4.02.3333",
			expected:    "2222222-33.2000.4.02.3333",
			expectError: false,
		},
		{
			name:        "Caso de borda - seg5=01 para seg4=4",
			input:       "3333333-44.2001.4.01.2222",
			expected:    "3333333-44.2001.4.01.2222",
			expectError: false,
		},
		{
			name:        "Caso de borda - seg5=05 para seg4=4",
			input:       "4444444-55.2001.4.05.1111",
			expected:    "4444444-55.2001.4.05.1111",
			expectError: false,
		},
		{
			name:        "Caso de borda - seg5=01 para seg4=8",
			input:       "5555555-66.2001.8.01.0000",
			expected:    "5555555-66.2001.8.01.0000",
			expectError: false,
		},
		{
			name:        "Caso de borda - seg5=27 para seg4=8",
			input:       "6666666-77.2001.8.27.9999",
			expected:    "6666666-77.2001.8.27.9999",
			expectError: false,
		},

		// Casos Inválidos
		{
			name:        "Formato inválido - tamanho errado (19 dígitos)",
			input:       "1234567-89.2001.4.03.567",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Formato inválido - caracteres não numéricos",
			input:       "12345a7-89.2001.4.03.5678",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Formato inválido - seg3 < 2000",
			input:       "1234567-89.1999.4.03.5678",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Formato inválido - seg4 diferente de '4' ou '8'",
			input:       "3333333-44.2005.7.01.2222",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Formato inválido - seg5 fora do intervalo para seg4=4 (seg5=06)",
			input:       "5555555-66.2022.4.06.0000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Formato inválido - seg5 fora do intervalo para seg4=8 (seg5=28)",
			input:       "4444444-55.2010.8.28.1111",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateFormat(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas obteve sucesso com resultado: %s", result)
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas obteve: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Resultado esperado: %s, mas obteve: %s", tt.expected, result)
				}
			}
		})
	}
}

// TestGenerateRandomPort testa a função GenerateRandomPort
func TestGenerateRandomPort(t *testing.T) {
	quantity := 5
	start := 7000 // Não está sendo usado na função atual, mas mantido para consistência

	ports := make(map[int]bool)
	for i := 0; i < quantity; i++ {
		port := GenerateRandomPort(quantity, start)
		if port < 7000 || port > 7999 {
			t.Errorf("Porta gerada fora do intervalo esperado: %d", port)
		}
		if ports[port] {
			t.Errorf("Porta duplicada gerada: %d", port)
		}
		ports[port] = true
	}
}

// TestExtractDriver testa a função ExtractDriver
func TestExtractDriver(t *testing.T) {
	driverPath, err := ExtractDriver()
	if err != nil {
		t.Fatalf("Falha ao extrair o driver: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(driverPath)) // Limpa o diretório temporário após o teste

	// Verifica se o arquivo existe
	if _, err := os.Stat(driverPath); os.IsNotExist(err) {
		t.Errorf("Driver não encontrado no caminho: %s", driverPath)
	}

	// Verifica permissões no Unix
	if runtime.GOOS != "windows" {
		info, err := os.Stat(driverPath)
		if err != nil {
			t.Errorf("Falha ao obter informações do arquivo: %v", err)
		}
		mode := info.Mode()
		if mode.Perm()&0755 != 0755 {
			t.Errorf("Permissões do arquivo incorretas: %v", mode.Perm())
		}
	}
}

// Mock Scraper e suas implementações para testar ScraperDispatcher

type MockPjeRJ struct{}

func (m *MockPjeRJ) Scrape(lawsuit string) (Hearing, error) {
	return Hearing{}, nil
}

type MockEsajSP struct{}

func (m *MockEsajSP) Scrape(lawsuit string) (Hearing, error) {
	return Hearing{}, nil
}

// TestScraperDispatcher testa a função ScraperDispatcher
func TestScraperDispatcher(t *testing.T) {
	tests := []struct {
		name         string
		lawsuit      string
		expectedType string // "PjeRJ", "EsajSP", "TJRJ", "Invalid"
		expectError  bool
	}{
		{
			name:         "Valid PJE-RJ format",
			lawsuit:      "08xxxxx-xx.20xx.8.19.xxxx", // Substitua x's por dígitos para corresponder à regex
			expectedType: "PjeRJ",
			expectError:  false,
		},
		{
			name:         "Valid ESAJ-SP format",
			lawsuit:      "0xxxxxx-xx.20xx.8.26.xxxx", // Substitua x's por dígitos para corresponder à regex
			expectedType: "EsajSP",
			expectError:  false,
		},
		{
			name:         "Valid TJRJ format",
			lawsuit:      "0[01]xxxxxxx20xx819xxxx", // Substitua [01] e x's por dígitos para corresponder à regex
			expectedType: "TJRJ",
			expectError:  false,
		},
		{
			name:         "Invalid format",
			lawsuit:      "invalid-format",
			expectedType: "Invalid",
			expectError:  true,
		},
		{
			name:         "Empty string",
			lawsuit:      "",
			expectedType: "Invalid",
			expectError:  true,
		},
		{
			name:         "Partial match",
			lawsuit:      "08xxxxx-xx.20xx.8.19",
			expectedType: "Invalid",
			expectError:  true,
		},
	}

	// Ajusta os valores de lawsuit para corresponder às regex
	for _, tt := range tests {
		switch tt.expectedType {
		case "PjeRJ":
			tt.lawsuit = "0800000-00.2000.8.19.0000"
		case "EsajSP":
			tt.lawsuit = "0123456-78.2001.8.26.1234"
		case "TJRJ":
			tt.lawsuit = "01234567820018191234"
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := ScraperDispatcher(tt.lawsuit)
			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas obteve scraper: %v", scraper)
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas obteve: %v", err)
				}
				switch tt.expectedType {
				case "PjeRJ":
					if _, ok := scraper.(*MockPjeRJ); !ok {
						t.Errorf("Esperava tipo *MockPjeRJ, mas obteve %T", scraper)
					}
				case "EsajSP":
					if _, ok := scraper.(*MockEsajSP); !ok {
						t.Errorf("Esperava tipo *MockEsajSP, mas obteve %T", scraper)
					}
				case "TJRJ":
					if scraper != nil {
						t.Errorf("Esperava scraper nulo para TJRJ, mas obteve: %T", scraper)
					}
				}
			}
		})
	}
}

// TestMatch testa a função match
func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected bool
	}{
		{
			name:     "Regex match success",
			pattern:  `^\d{7}\d{2}\d{4}8{1}19{1}\d{4}$`,
			input:    "1234567891238195678",
			expected: true,
		},
		{
			name:     "Regex match failure",
			pattern:  `^\d{7}\d{2}\d{4}8{1}19{1}\d{4}$`,
			input:    "12345678A1238195678",
			expected: false,
		},
		{
			name:     "Invalid regex pattern",
			pattern:  `^\d{7}(`, // Regex inválida
			input:    "1234567891238195678",
			expected: false,
		},
		{
			name:     "Empty input",
			pattern:  `^\d+$`,
			input:    "",
			expected: false,
		},
		{
			name:     "Empty pattern",
			pattern:  ``,
			input:    "any input",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Match(tt.pattern, tt.input)
			if result != tt.expected {
				t.Errorf("Esperava %v, mas obteve %v para input '%s' e pattern '%s'", tt.expected, result, tt.input, tt.pattern)
			}
		})
	}
}
