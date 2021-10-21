package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
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

	// open db connection
	db, err := sql.Open("mysql", "root:pa55w0rd@tcp(localhost:3306)/test_db")
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	//err = db.Ping()
	//if err != nil {
	//	panic(err.Error())
	//}

	// exec query
	rows, err := db.Query("select * from content")
	if err != nil {
		panic(err.Error())
	}

	// get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	//make slice for the values
	values := make([]sql.RawBytes, len(columns))

	// row.Scan wants []inteface{} as an argument, so we must copy the references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// fetch rows
	for rows.Next() {
		// get raw bytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// now do smth with data
		// print, for example
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("--------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

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
