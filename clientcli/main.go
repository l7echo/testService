package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
)

func main() {
	var restParams restAction
	err := getInputParams(&restParams)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		flag.Usage()
		os.Exit(1)
	}

	var client Client
	client.BaseURL = "http://" + restParams.host + ":" + restParams.port
	client.HTTPClient = http.DefaultClient

	reflectiveResult, err := call(restParams.actionName, &client, &restParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	result := reflectiveResult[0].Interface().([]Record)

	if reflectiveResult[1].Interface() != nil {
		errRet := reflectiveResult[1].Interface().(error)
		fmt.Printf("error: %s\n", errRet)
	}
	fmt.Printf("result: %s\n", result)
}

//---------------------------------------------------------
// input params handler

type restAction struct {
	host       string
	port       string
	id         string
	value      string
	actionName string
}

func getInputParams(params *restAction) error {

	// maybe using smth like urfave/cli would be better choice
	// but it is too big dependency for this short program

	var actionAdd bool
	var actionGet bool
	var actionGetAll bool
	var actionRemove bool

	var (
		host  = flag.String("host", "localhost", "")
		port  = flag.String("port", "19300", "")
		id    = flag.String("id", "", "")
		value = flag.String("value", "", "")
	)
	flag.BoolVar(&actionAdd, "add", false, "")
	flag.BoolVar(&actionGet, "get", false, "")
	flag.BoolVar(&actionGetAll, "get-all", false, "")
	flag.BoolVar(&actionRemove, "remove", false, "")

	flag.Usage = func() {
		fmt.Printf("Usage: \n")
		fmt.Printf("   %s --host=... --port=... --id=... --value=... --add/get/get-all/remove\n", os.Args[0])
		fmt.Printf("Where: \n")
		fmt.Printf("   --host using to set backend host address. Default value is localhost\n")
		fmt.Printf("   --port using to set backend port number. Default value is 19300\n")
		fmt.Printf("   --id using for ID of DB content. Requed by add, get, remove\n")
		fmt.Printf("   --value using for DB content. Requed by add, get, remove\n")
		fmt.Printf("   --add/get/get-all/remove set your rest action\n")
	}
	flag.Parse()

	actionFromInput := make(map[string]bool)
	actionFromInput["add"] = actionAdd
	actionFromInput["get"] = actionGet
	actionFromInput["get-all"] = actionGetAll
	actionFromInput["remove"] = actionRemove

	var selectedAction string
	for i, v := range actionFromInput {
		if v == true {
			if selectedAction != "" {
				return errors.New("you have an error in params! Please, select ONE action\n")
			}
			selectedAction = i
		}
	}
	if selectedAction == "" {
		return errors.New("you must set action param\n")
	}

	// here i need to check id & value
	// but more useful will be doing this check later: in action function

	params.host = *host
	params.port = *port
	params.id = *id
	params.value = *value
	params.actionName = selectedAction

	return nil
}

//---------------------------------------------------------
// client rest api

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Record struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

var restActionsMap = map[string]interface{}{
	"get":     (*Client).get,
	"get-all": (*Client).getAll,
}

func (client *Client) get(restParam *restAction) ([]Record, error) {
	var req *http.Request
	var err error

	// if ID and VALUE unset both
	if restParam.id == "" && restParam.value == "" {
		err = errors.New("get: you need to set ID and/or VALUE\n")
		return nil, err
	}

	// if ID and VALUE set both
	if restParam.id != "" && restParam.value != "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/get/%s/%s", client.BaseURL, restParam.id, restParam.value), nil)
		if err != nil {
			return nil, err
		}
	}

	if restParam.id != "" && restParam.value == "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/get/id=%s", client.BaseURL, restParam.id), nil)
		if err != nil {
			return nil, err
		}
	}

	if restParam.id == "" && restParam.value != "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/get/value=%s", client.BaseURL, restParam.value), nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("check backend, status code: %d", res.StatusCode)
	}
	var resList []Record
	var data Record
	err = json.NewDecoder(res.Body).Decode(&data)
	resList = append(resList, data)

	return resList, err
}

func (client *Client) getAll(restParam *restAction) ([]Record, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/get-all", client.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("check backend, status code: %d", res.StatusCode)
	}
	var resList []Record
	err = json.NewDecoder(res.Body).Decode(&resList)

	return resList, nil
}

// the main idea: we detect needed function by name and call it
// needed function name set by input argument
// function arguments count may be variadic
func call(name string, params ...interface{}) (result []reflect.Value, err error) {
	function := reflect.ValueOf(restActionsMap[name])

	if len(params) != function.Type().NumIn() {
		err = errors.New("call(): wrong parameters count\n")
		return nil, err
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = function.Call(in)

	return result, nil
}
