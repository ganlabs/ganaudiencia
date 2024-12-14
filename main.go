package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

func getHearing(w http.ResponseWriter, r *http.Request) {

	lawsuit := r.URL.Query().Get("processo")
	if lawsuit == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l, err := ValidateFormat(lawsuit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	mv, err := FetchMovements(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	hd, ht := ExtractHearingDates(mv)

	log.Println(hd, ht)

	h := Hearing{
		Lawsuit:     lawsuit,
		Class:       ExtractClass(mv[0]),
		HearingDate: hd,
		HearingTime: ht,
		IsValid:     ValidateDate(hd),
		Movement:    mv,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h)

}

//go:embed static/*
var staticFiles embed.FS

func main() {

	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(staticContent))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /audiencia", getHearing)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// read the template from a file or a string
		t, err := template.ParseFS(staticFiles, "static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	if err := http.ListenAndServe(":3003", LoggingMiddleware(mux)); err != nil {
		panic(err)
	}

}
