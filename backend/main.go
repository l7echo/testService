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

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {

		// get by id
		if item.Id == vars["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}

		// get by value
		if item.Value == vars["value"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// get-all case
	json.NewEncoder(w).Encode(dbContent)
}

func main() {
	r := mux.NewRouter()

	dbContent = append(dbContent, Record{Id: "1", Value: "value1"})
	dbContent = append(dbContent, Record{Id: "2", Value: "value2"})

	r.HandleFunc("/get/id={id}", getHandler).Methods("GET")
	r.HandleFunc("/get/value={value}", getHandler).Methods("GET")
	r.HandleFunc("/get-all", getHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":19300", r))
}
