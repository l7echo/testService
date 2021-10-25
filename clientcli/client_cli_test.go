package main

import (
	"net/http"
	"os"
	"testing"
)

func getClient(params *restAction) *Client {
	var client Client
	client.BaseURL = "http://" + params.host + ":" + params.port
	client.HTTPClient = http.DefaultClient
	return &client
}

func getTestParams(params *restAction) {
	params.host = "localhost"
	params.port = "19300"
	params.id = "test_id"
	params.value = "test_value"
}

func TestAdd(t *testing.T) {
	var testParams restAction
	getTestParams(&testParams)

	client := getClient(&testParams)

	if ok := client.checkServer(&testParams); ok != true {
		t.Error("Server is not avaible\n")
		os.Exit(1)
	}

	var testRecord Record
	testRecord.Id = testParams.id
	testRecord.Value = testParams.value

	rec, err := client.add(&testParams)

	if err != nil {

	}

	if len(rec) != 1 {

	}
}

func TestGet(t *testing.T) {

}

func TestGetAll(t *testing.T) {

}

func TestRemove(t *testing.T) {

}
