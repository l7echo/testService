package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// read db config & open db pool connection
func openDB() (*sql.DB, error) {
	dbConfig, err := readDBParams()
	if err != nil {
		return nil, err
	}
	dbConnPool, err = openDbPool(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.DbName))
	if err != nil {
		return nil, err
	}
	return dbConnPool, nil
}

func isSubStr(str, substr string) bool {
	for i := 0; i < len(str); i++ {
		if strings.HasPrefix(str[i:], substr) {
			return true
		}
	}
	return false
}

func TestDbConnection(t *testing.T) {
	var originArg1 string
	originArg1 = os.Args[1]
	var needToRestoreArgs bool

	if isSubStr(os.Args[1], "dbconfig.json") == false {
		t.Log("didn't set \"./config/dbconfig.json\" parameter, let's fix it\n")
		os.Args[1] = "./config/dbconfig.json"
		needToRestoreArgs = true
	}

	dbConnPool, err := openDB()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = dbConnPool.Close()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if needToRestoreArgs {
		os.Args[1] = originArg1
	}
}

func TestGetRequests(t *testing.T) {
	var originArg1 string
	originArg1 = os.Args[1]
	var needToRestoreArgs bool

	if isSubStr(os.Args[1], "dbconfig.json") == false {
		t.Log("didn't set \"./config/dbconfig.json\" parameter, let's fix it\n")
		os.Args[1] = "./config/dbconfig.json"
		needToRestoreArgs = true
	}

	dbConnPool, err := openDB()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	defer dbConnPool.Close()

	routerVars := []struct {
		variable   string
		shouldPass bool
		f          func(http.ResponseWriter, *http.Request)
	}{
		// right data
		{"/db/id=1", true, getById},
		{"/db/value=v2", true, getByValue},
		{"/db/3/v3", true, getByIdAndValue},
		{"/db", true, getAll},
		// wrong data
		{"/bdbdbd", false, getAll},
	}

	for _, routerVariable := range routerVars {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s", routerVariable.variable), nil)
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		respRecoder := httptest.NewRecorder()

		// we need to create a router because we need to add our requests into context
		router := mux.NewRouter()
		router.HandleFunc("/db/id={id}", getById).Methods("GET")
		router.HandleFunc("/db/value={value}", getByValue).Methods("GET")
		router.HandleFunc("/db/{id}/{value}", getByIdAndValue).Methods("GET")
		router.HandleFunc("/db", getAll).Methods("GET")
		router.ServeHTTP(respRecoder, req)

		if status := respRecoder.Code; status != http.StatusOK && !routerVariable.shouldPass {

			// check if it is a special wrong data
			// if it is, we passed this test
			if routerVariable.variable != "/bdbdbd" {
				t.Errorf("got wrong http status: got %d, want %d", status, http.StatusOK)
				t.FailNow()
			}
		}
	}

	if needToRestoreArgs {
		os.Args[1] = originArg1
	}
}

func TestPostRequest(t *testing.T) {
	var originArg1 string
	originArg1 = os.Args[1]
	var needToRestoreArgs bool

	if isSubStr(os.Args[1], "dbconfig.json") == false {
		t.Log("didn't set \"./config/dbconfig.json\" parameter, let's fix it\n")
		os.Args[1] = "./config/dbconfig.json"
		needToRestoreArgs = true
	}

	dbConnPool, err := openDB()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	defer dbConnPool.Close()

	records := []struct {
		record Record
		status int
	}{
		{Record{Id: "test_id", Value: "test_value"}, http.StatusOK},
		{Record{Id: "test_id", Value: "test_value"}, http.StatusInternalServerError},
	}

	for i := 0; i < len(records); i++ {
		jsonStr, _ := json.Marshal(records[i].record)

		req, err := http.NewRequest("POST", "/db", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		respRecoder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/db", addNew).Methods("POST")
		router.ServeHTTP(respRecoder, req)

		if status := respRecoder.Code; status != records[i].status {
			t.Errorf("got wrong http status: got %d, want %d", status, records[i].status)
			t.FailNow()
		}
	}

	if needToRestoreArgs {
		os.Args[1] = originArg1
	}
}

func TestDeleteRequest(t *testing.T) {
	var originArg1 string
	originArg1 = os.Args[1]
	var needToRestoreArgs bool

	if isSubStr(os.Args[1], "dbconfig.json") == false {
		t.Log("didn't set \"./config/dbconfig.json\" parameter, let's fix it\n")
		os.Args[1] = "./config/dbconfig.json"
		needToRestoreArgs = true
	}

	dbConnPool, err := openDB()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	defer dbConnPool.Close()

	records := []struct {
		recordId string
		status   int
	}{
		{"test_id", http.StatusOK},
		//{"test_id", http.StatusInternalServerError},
	}

	for i := 0; i < len(records); i++ {
		// let's delete Record{Id: "test_id", Value: "test_value"} from previous test
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/db/%s", records[i].recordId), nil)
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		respRecoder := httptest.NewRecorder()

		// we need to create a router because we need to add our requests into context
		router := mux.NewRouter()
		router.HandleFunc("/db/{id}", deleteById).Methods("DELETE")
		router.ServeHTTP(respRecoder, req)

		if status := respRecoder.Code; status != records[i].status {
			t.Errorf("got wrong http status: got %d, want %d", status, records[i].status)
			t.FailNow()
		}
	}

	if needToRestoreArgs {
		os.Args[1] = originArg1
	}
}
