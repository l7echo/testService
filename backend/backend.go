package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

var dbConnPool *sql.DB

func getById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	rows, err := dbConnPool.Query("call searchById(?)", vars["id"]) // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close() // return connection to pool

	dbContentLocal, err := scanDbOutput(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(*dbContentLocal)
}

// we can get a lot of non-unique values from db
func getByValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	rows, err := dbConnPool.Query("call searchByValue(?)", vars["value"]) // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close() // return connection to pool

	dbContentLocal, err := scanDbOutput(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(*dbContentLocal)
}

func getByIdAndValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	row, err := dbConnPool.Query("call searchByIdAndValue(?, ?)", vars["id"], vars["value"]) // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer row.Close() // return connection to pool

	dbContentLocal, err := scanDbOutput(row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(*dbContentLocal)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := dbConnPool.Query("call getAll()") // call store procedure
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close() // return connection to pool

	dbContentLocal, err := scanDbOutput(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(*dbContentLocal)
}

func addNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var record Record
	_ = json.NewDecoder(r.Body).Decode(&record)

	row, err := dbConnPool.Query("call addValue(?, ?)", record.Id, record.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	err = row.Close() // return connection to pool
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(record) // return added record
}

func deleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	row, err := dbConnPool.Query("call deleteById(?)", vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	err = row.Close() // return connection to pool
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(&Record{}) // return empty record
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
	file, err := ioutil.ReadFile(os.Args[1])
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

func scanDbOutput(rows *sql.Rows) (*[]Record, error) {
	var record Record
	var dbContentLocal []Record

	for rows.Next() {
		err := rows.Scan(&record.Id, &record.Value)
		if err != nil {
			return nil, err
		} else {
			if record.Id == "" && record.Value == "" {
				dbContentLocal = append(dbContentLocal, Record{})
			} else {
				dbContentLocal = append(dbContentLocal, record)
			}
		}
	}

	return &dbContentLocal, nil
}

func main() {
	if len(os.Args) != 2 {
		panic(errors.New("you must set json with DB connection parameters as argument\n"))
	}
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

	r.HandleFunc("/db/id={id}", getById).Methods("GET")
	r.HandleFunc("/db/value={value}", getByValue).Methods("GET")
	r.HandleFunc("/db/{id}/{value}", getByIdAndValue).Methods("GET")
	r.HandleFunc("/db", getAll).Methods("GET")
	r.HandleFunc("/db", addNew).Methods("POST")
	r.HandleFunc("/db/{id}", deleteById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":19300", r))
}
