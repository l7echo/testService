package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

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
		fmt.Printf("   %s --host=... --port=... [--id=...] [--value=...] --add/get/get-all/remove\n", os.Args[0])
		fmt.Printf("Where: \n")
		fmt.Printf("   --host using to set backend host address. Default value is localhost\n")
		fmt.Printf("   --port using to set backend port number. Default value is 19300\n")
		fmt.Printf("   --id using for ID of DB content. Can be requed by last parameter\n")
		fmt.Printf("   --value using for DB content. Can be requed by last parameter\n")
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

	params.host = *host
	params.port = *port
	params.id = *id
	params.value = *value
	params.actionName = selectedAction

	return nil
}

func main() {
	var restParams restAction
	err := getInputParams(&restParams)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		flag.Usage()
		os.Exit(1)
	}
}
