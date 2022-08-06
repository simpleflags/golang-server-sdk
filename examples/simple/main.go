package main

import (
	"context"
	"fmt"
	"github.com/simpleflags/golang-server-sdk/client"
	"github.com/simpleflags/golang-server-sdk/connector/simple"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const sdkKey = "12d466a8-f62a-11ec-b4e5-faffc22119b1"

const featureFlagKey = "bool-flag"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	target := map[string]interface{}{
		"identifier": "enver",
	}

	conn := simple.NewHttpConnector(sdkKey, simple.WithStreamURL("https://64a55c46.fanoutcdn.com/api"))
	sf, err := client.NewWithConnector(conn)

	if err != nil {
		log.Printf("could not connect to SF servers %v", err)
	}

	defer func() {
		if err := sf.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()
	sf.WaitForInitialization()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				showFeature := sf.Evaluate(featureFlagKey, target).Bool(false)

				fmt.Printf("KeyFeature flag '%s' is %t for this user\n", featureFlagKey, showFeature)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	<-ctx.Done()
}
