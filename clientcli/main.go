package main

import (
	"errors"
	"flag"
	"fmt"
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

	result, err := call(restParams.actionName, &restParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
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
				return errors.New("you have an error in last param! Please, select ONE action\n")
			}
			selectedAction = i
		}
	}
	if selectedAction == "" {
		return errors.New("you must set last param\n")
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

var restActionsMap = map[string]interface{}{
	"get":     get,
	"get-all": getAll,
}

func get(restParam *restAction) {

}

func getAll(restParam *restAction) {

}

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
