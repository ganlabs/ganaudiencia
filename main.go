package main

import (
	"log"
	"net/http"
)

func getHearing(w http.ResponseWriter, r *http.Request) {

	lawsuit := r.URL.Query().Get("audiencia")
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(l))

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", getHearing)
	if err := http.ListenAndServe(":3003", LoggingMiddleware(mux)); err != nil {
		panic(err)
	}

}
