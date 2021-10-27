package main

import (
	"flag"
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
		t.FailNow()
	}

	return &testParams, client
}

func TestCall(t *testing.T) {
	testParams, client := createTestData(t)

	// at the first i would like to check call() func
	// and prepare DB content for all other tests a little

	reflectiveResult, err := call("remove", client, testParams)
	if err != nil {
		t.Errorf("error: %s\n", err)
		t.FailNow()
	}

	if reflectiveResult[1].Interface() != nil {
		errRet := reflectiveResult[1].Interface().(error)
		t.Errorf("error: %s\n", errRet)
		t.Fail()
	} else {
		result := reflectiveResult[0].Interface().([]Record)
		emptyRecord := Record{Id: "", Value: ""}
		if result[0] != emptyRecord {
			t.Error("call(client.remove) failed\n")
			t.Fail()
		}
	}
}

func TestAdd(t *testing.T) {
	testParams, client := createTestData(t)
	rec, err := client.add(testParams)

	if err != nil {
		t.Errorf("client.add failed: %s\n", err)
		t.Fail()
		return
	}

	if len(rec) != 1 {
		t.Errorf("client.add failed: len(rec)=%d\n", len(rec))
		t.Fail()
		return
	}

	testParams.value = ""
	_, err = client.add(testParams) // and let's check invalid input data

	if err == nil {
		t.Error("client.add(id=nil) failed\n")
		t.Fail()
	}
}

func isSubStr(str, substr string) bool {
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
		if isSubStr(err.Error(), "Duplicate entry 'test_id' for key 'content.PRIMARY'") == false {
			t.Errorf("client.add return: %s\n", err.Error())
			t.Fail()
		}
	} else {
		t.Error("client.add failed: addition duplicate data was success\n")
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	testParams, client := createTestData(t)

	var testRecord Record
	testRecord.Id = testParams.id
	testRecord.Value = testParams.value

	rec, err := client.get(testParams) // here we have id & value, so we will get 1 record

	if err != nil {
		t.Errorf("client.get(id, value) failed: %s\n", err)
		t.Fail()
		return
	}

	if len(rec) != 1 {
		t.Errorf("client.get(id, value) failed: len(rec)=%d\n", len(rec))
		t.Fail()
		return
	}

	if rec[0] != testRecord {
		t.Error("client.get(id, value) failed: test data is not equal DB return\n")
		t.Fail()
		return
	}

	testParams.value = ""
	rec, err = client.get(testParams) // here we have id, so we will get 1 record by unique id

	if err != nil {
		t.Errorf("client.get(id) failed: %s\n", err)
		t.Fail()
		return
	}

	if len(rec) != 1 {
		t.Errorf("client.get(id) failed: len(rec)=%d\n", len(rec))
		t.Fail()
		return
	}

	if rec[0] != testRecord {
		t.Error("client.get(id) failed: test data is not equal DB return\n")
		t.Fail()
		return
	}

	testParams.value = "test_value"
	testParams.id = ""
	rec, err = client.get(testParams) // here we have only a value, so we may get several records

	if err != nil {
		t.Errorf("client.get(value) failed: %s\n", err)
		t.Fail()
		return
	}
	var isTestPass bool
	for _, record := range rec {
		if record == testRecord {
			isTestPass = true
		}
	}
	if isTestPass == false {
		t.Error("client.get(value) failed: test data is not in DB return\n")
		t.Fail()
		return
	}

	testParams.value = ""
	testParams.id = ""
	rec, err = client.get(testParams) // and let's check invalid input data

	if err == nil || rec != nil {
		t.Error("client.get(id=nil, value=nil) failed\n")
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	testParams, client := createTestData(t)

	var testRecord Record
	testRecord.Id = testParams.id
	testRecord.Value = testParams.value

	rec, err := client.getAll(testParams) // here we may get several records

	if err != nil {
		t.Errorf("client.remove failed, err: %s\n", err.Error())
		t.Fail()
		return
	}

	var isTestPass bool
	for _, record := range rec {
		if record == testRecord {
			isTestPass = true
		}
	}
	if isTestPass == false {
		t.Error("client.get(value) failed: test data is not in DB return\n")
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	testParams, client := createTestData(t)

	rec, err := client.remove(testParams)

	if err != nil {
		t.Errorf("client.remove failed, err: %s\n", err.Error())
		t.Fail()
		return
	}
	if len(rec) != 1 {
		t.Error("client.remove failed\n")
		t.Fail()
		return
	}
	emptyRecord := Record{Id: "", Value: ""}
	if rec[0] != emptyRecord {
		t.Error("client.remove failed\n")
		t.Fail()
		return
	}

	testParams.id = ""
	rec, err = client.remove(testParams) // and let's check invalid input data

	if err == nil || rec != nil {
		t.Error("client.remove(id=nil) failed\n")
		t.Fail()
	}
}

func TestGetBadInputParams(t *testing.T) {
	var inputParams restAction

	// may be redefinition os.Args is a very bad idea
	originArgs := make([]string, len(os.Args))
	copy(originArgs, os.Args)

	os.Args = os.Args[:1]

	// set wrong input
	os.Args = append(os.Args, "--get")
	os.Args = append(os.Args, "--add")
	err := getInputParams(&inputParams) // and now check it

	if err.Error() != "you have an error in params! Please, select ONE action\n" {
		os.Args = os.Args[:0]
		copy(os.Args, originArgs)

		t.Error("getInputParams failed with wrong input (with many actions)\n")
		t.Fail()
		return
	}

	// restore
	os.Args = os.Args[:0]
	copy(os.Args, originArgs)
}

func TestCallingMain(tester *testing.T) {
	originArgs := make([]string, len(os.Args))
	copy(originArgs, os.Args)

	os.Args = os.Args[:1]

	// set input
	os.Args = append(os.Args, "--get-all")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // we need to clear flag after TestGetBadInputParams

	main()

	// restore
	os.Args = os.Args[:0]
	copy(os.Args, originArgs)
}
