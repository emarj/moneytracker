package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"ronche.se/moneytracker/model"
)

func JSONHandler(srv model.Service) http.Handler {
	mux := http.NewServeMux()
	//h := jsonHandler{srv}
	return mux
}

type jsonHandler struct {
	srv model.Service
	//Google Sheets
}

func jsonResponseWriter(f func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, status, err := f(w, r)
		if err != nil {
			log.Printf("Error (%d) %v", status, err)
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}

	}
}
