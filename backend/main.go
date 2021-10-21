package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Record struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type DBconfig struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	DbName string `json:"dbname"`
}

var dbContent []Record

var dbConnPool *sql.DB

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

// we can get a lot of non-unique values from db
func getByValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var dbContentLocal []Record

	rows, err := dbConnPool.Query("call searchByValue(?)", vars["value"]) // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	columns, err := rows.Columns() // columns names
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	// row.Scan wants []inteface{} as an argument, so we must copy the references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get raw bytes from data
		err = rows.Scan(scanArgs...)
		var record Record

		err = rows.Scan(&record.Id, &record.Value)
		if err != nil {
			panic(err.Error())
		}

		// now do save data from sql db to local record var
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if columns[i] == "id" {
				record.Id = value
			}
			if columns[i] == "value" {
				record.Value = value
			}
		}
		dbContentLocal = append(dbContentLocal, record)
	}

	_ = json.NewEncoder(w).Encode(dbContentLocal)
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

	var dbContentLocal []Record

	rows, err := dbConnPool.Query("call getAll()") // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}

	for rows.Next() {
		var record Record

		err = rows.Scan(&record.Id, &record.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err.Error())
		} else {
			dbContentLocal = append(dbContentLocal, record)
		}
	}

	_ = json.NewEncoder(w).Encode(dbContentLocal)
}

func addNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var record Record
	_ = json.NewDecoder(r.Body).Decode(&record)

	_, err := dbConnPool.Query("call addValue(?, ?)", record.Id, record.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}
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

func openDbPool(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	// let's configure our pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	return db, nil
}

func readDBParams() (*DBconfig, error) {
	file, err := ioutil.ReadFile("C:\\Projects\\go\\src\\testService\\backend\\dbconfig.json")
	if err != nil {
		return nil, err
	}

	data := DBconfig{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func main() {
	// read db config & open db pool connection
	dbConfig, err := readDBParams()
	if err != nil {
		panic(err.Error())
	}

	dbConnPool, err = openDbPool(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.DbName))
	if err != nil {
		panic(err.Error())
	}
	defer dbConnPool.Close()

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
