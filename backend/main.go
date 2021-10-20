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

func getById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Id == vars["id"] {
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&Record{})
}

func getByValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Value == vars["value"] {
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&Record{})
}

func getByIdAndValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, item := range dbContent {
		if item.Id == vars["id"] && item.Value == vars["value"] {
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&Record{})
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dbContent)
}

func addNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var record Record
	_ = json.NewDecoder(r.Body).Decode(&record)
	dbContent = append(dbContent, record)
	_ = json.NewEncoder(w).Encode(record)
}

func deleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for index, item := range dbContent {
		if item.Id == vars["id"] {
			dbContent = append(dbContent[:index], dbContent[index+1:]...)
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&Record{})
}

func main() {
	r := mux.NewRouter()

	dbContent = append(dbContent, Record{Id: "1", Value: "v1"})
	dbContent = append(dbContent, Record{Id: "2", Value: "v2"})

	r.HandleFunc("/db/id={id}", getById).Methods("GET")
	r.HandleFunc("/db/value={value}", getByValue).Methods("GET")
	r.HandleFunc("/db/{id}/{value}", getByIdAndValue).Methods("GET")
	r.HandleFunc("/db", getAll).Methods("GET")
	r.HandleFunc("/db", addNew).Methods("POST")
	r.HandleFunc("/db/{id}", deleteById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":19300", r))
}
