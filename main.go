package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

type cacheEntry struct {
	hearing  Hearing
	cachedAt time.Time
}

var (
	cache       = make(map[string]*cacheEntry)
	cacheMu     sync.RWMutex
	Environment = os.Getenv("ENVIRONMENT")
	BaseURL     = os.Getenv("BASE_URL")
	port        int
)

//go:embed driver/*
var Chromedriver embed.FS

const cacheTTL = 3 * time.Hour

func getHearing(w http.ResponseWriter, r *http.Request) {
	lawsuit := r.URL.Query().Get("processo")
	if lawsuit == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	noCache := r.URL.Query().Get("nocache") == "true"

	if !noCache {
		cacheMu.RLock()
		ce, found := cache[lawsuit]
		cacheMu.RUnlock()

		if found {
			if time.Since(ce.cachedAt) <= cacheTTL {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(ce.hearing)
				return
			} else {
				cacheMu.Lock()
				delete(cache, lawsuit)
				cacheMu.Unlock()
			}
		}
	}

	l, err := ValidateFormat(lawsuit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	sc, err := ScraperDispatcher(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	h, err := sc.Scrape(l)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !noCache {
		cacheMu.Lock()
		cache[lawsuit] = &cacheEntry{
			hearing:  h,
			cachedAt: time.Now(),
		}
		cacheMu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h)
}

//go:embed static/*
var staticFiles embed.FS

func main() {

	corsConfig := CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3000", "http://localhost:5173", "https://seusite.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600, // 1 hora
	}

	switch Environment {
	case "development":
		log.Println("Executando em ambiente de desenvolvimento.")
		port = 3003
		BaseURL = fmt.Sprintf("http://localhost:%d", port)
	case "docker":
		log.Println("Executando em ambiente docker.")
		port = 3003
	default:
		log.Println("Executando em ambiente local.")
		port = GenerateRandomPort(998, 7001)
		BaseURL = fmt.Sprintf("http://localhost:%d", port)
	}

	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticContent))

	mux := http.NewServeMux()

	mux.HandleFunc("/audiencia", getHearing)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(staticFiles, "static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"BaseURL": BaseURL,
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.Handle("/sair", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.instagram.com/gondim_adv/", http.StatusFound)

		if os.Getenv("ENVIRONMENT") != "docker" {
			go func() {
				time.Sleep(1 * time.Second)
				fmt.Println("Shutting down the backend.")
				os.Exit(0)
			}()
		}
	}))

	loggedHandler := LoggingMiddleware(CORSMiddleware(corsConfig, mux))

	if os.Getenv("ENVIRONMENT") != "docker" {
		go (func() {
			time.Sleep(1 * time.Second)
			if err := OpenBrowser(fmt.Sprintf("http://localhost:%d", port)); err != nil {
				log.Println(err)
			}
		})()
	}

	log.Printf("Servidor iniciado na porta %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), loggedHandler); err != nil {
		panic(err)
	}
}
