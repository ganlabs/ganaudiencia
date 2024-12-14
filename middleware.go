package main

import (
	"bytes"
	"encoding/json"
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
			"method":        r.Method,
			"url":           r.URL.String(),
			"headers":       r.Header,
			"request_body":  requestBody,
			"status":        responseRecorder.statusCode,
			"response_body": responseRecorder.body.String(),
			"duration":      time.Since(start).String(),
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
