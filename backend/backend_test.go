package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
		handler := http.HandlerFunc(routerVariable.f)

		handler.ServeHTTP(respRecoder, req)
		if status := respRecoder.Code; status != http.StatusOK {
			t.Errorf("got wrong http status: got %d, want %d", status, http.StatusOK)
			t.FailNow()
		}
	}

	if needToRestoreArgs {
		os.Args[1] = originArg1
	}
}

func TestPostRequest(t *testing.T) {

}

func TestDeleteRequest(t *testing.T) {

}
