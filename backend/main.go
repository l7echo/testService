package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Record struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

var dbContent []Record

func getHandlerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Id == vars["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Record{})
}

func getHandlerByValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Value == vars["value"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Record{})
}

func getHandlerByIdAndValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Id == vars["id"] && item.Value == vars["value"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Record{})
}

func getHandlerAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbContent)
}

func main() {
	r := mux.NewRouter()

	dbContent = append(dbContent, Record{Id: "1", Value: "value1"})
	dbContent = append(dbContent, Record{Id: "2", Value: "value2"})

	r.HandleFunc("/get/id={id}", getHandlerById).Methods("GET")
	r.HandleFunc("/get/value={value}", getHandlerByValue).Methods("GET")
	r.HandleFunc("/get/{id}/{value}", getHandlerByIdAndValue).Methods("GET")
	r.HandleFunc("/get-all", getHandlerAll).Methods("GET")

	log.Fatal(http.ListenAndServe(":19300", r))
}
