package main

import (
	"net/http"
	"os"
	"strings"
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

func createTestData(t *testing.T) (*restAction, *Client) {
	var testParams restAction
	getTestParams(&testParams)

	client := getClient(&testParams)

	if ok := client.checkServer(&testParams); ok != true {
		t.Error("Server is not avaible\n")
		os.Exit(1)
	}

	return &testParams, client
}

func TestAdd(t *testing.T) {
	testParams, client := createTestData(t)
	rec, err := client.add(testParams)

	if err != nil {
		t.Errorf("client.add failed: %s", err)
		t.Fail()
		return
	}

	if len(rec) != 1 {
		t.Errorf("client.add failed: len(rec)=%d", len(rec))
		t.Fail()
		return
	}
}

func strstr(str, substr string) bool {
	for i := 0; i < len(str); i++ {
		if strings.HasPrefix(str[i:], substr) {
			return true
		}
	}
	return false
}

func TestDoubleAdd(t *testing.T) {
	testParams, client := createTestData(t)
	_, err := client.add(testParams)

	if err != nil {
		if strstr(err.Error(), "Duplicate entry 'test_id' for key 'content.PRIMARY'") == false {
			t.Errorf("client.add return: %s", err.Error())
			t.Fail()
		}
	} else {
		t.Error("client.add failed: addition duplicate data was success")
		t.Fail()
	}
}

func TestGet(t *testing.T) {

}

func TestGetAll(t *testing.T) {

}

func TestRemove(t *testing.T) {

}
