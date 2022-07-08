package main

import (
	"fmt"
	"github.com/simpleflags/golang-server-sdk/client"
	"github.com/simpleflags/golang-server-sdk/connector"
	"log"
)

const sdkKey = "12d466a8-f62a-11ec-b4e5-faffc22119b1"

const featureFlagKey = "bool-flag"

func main() {
	target := map[string]interface{}{
		"identifier": "enver",
	}

	conn, err := connector.NewFileConnector(sdkKey, "/Users/enver/Projects/sf-workspace/golang-server-sdk/examples/data")
	if err != nil {
		log.Println(err)
	}
	sf, err := client.NewWithConnector(conn, client.WithPullerEnabled(false), client.WithStreamEnabled(false),
		client.WithPrefetchFlags("bool-flag"))

	if err != nil {
		log.Printf("could not connect to SF servers %v", err)
	}

	defer func() {
		if err := sf.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()

	showFeature := sf.Evaluate(featureFlagKey, target).Bool(false)

	fmt.Printf("KeyFeature flag '%s' is %t for this user\n", featureFlagKey, showFeature)
}
