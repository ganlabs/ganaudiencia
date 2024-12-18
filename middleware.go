package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var requestBody string
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for further use
				requestBody = string(bodyBytes)
			}
		}

		responseRecorder := &responseWriter{ResponseWriter: w, body: &bytes.Buffer{}, statusCode: http.StatusOK}
		next.ServeHTTP(responseRecorder, r)

		logEntry := map[string]interface{}{
			"method":       r.Method,
			"url":          r.URL.String(),
			"headers":      r.Header,
			"request_body": requestBody,
			"status":       responseRecorder.statusCode,
			// "response_body": responseRecorder.body.String(),
			"duration": time.Since(start).String(),
		}

		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
		} else {
			log.Println(string(logData))
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body.Write(body)
	return rw.ResponseWriter.Write(body)
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func CORSMiddleware(config CORSConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Verifica se a origem está na lista de permitidas
		allowed := false
		for _, o := range config.AllowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Vary", "Origin")

			if r.Method == http.MethodOptions {
				// Trata requisições pré-voo
				w.Header().Set("Access-Control-Allow-Methods", join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", join(config.AllowedHeaders, ", "))
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		// Passa para o próximo handler
		next.ServeHTTP(w, r)
	})
}

// join é uma função auxiliar para concatenar strings com um separador
func join(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}
