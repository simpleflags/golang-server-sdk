package main

import (
	"context"
	"fmt"
	sdk "github.com/simpleflags/golang-server-sdk"
	"github.com/simpleflags/golang-server-sdk/client"
	"github.com/simpleflags/golang-server-sdk/connector/simple"
	"github.com/simpleflags/golang-server-sdk/repository"
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

	fileStorage, err := repository.NewFileStorage("./")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	conn := simple.NewHttpConnector(sdkKey, simple.WithStreamURL("https://64a55c46.fanoutcdn.com/api"))
	err = sdk.InitWithConnector(conn, client.WithStorage(&fileStorage))
	if err != nil {
		log.Printf("could not connect to SF servers %v", err)
	}

	defer func() {
		if err := sdk.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()
	sdk.WaitForInitialization()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				showFeature := sdk.Evaluate(featureFlagKey, target).Bool(false)

				fmt.Printf("KeyFeature flag '%s' is %t for this user\n", featureFlagKey, showFeature)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	<-ctx.Done()
}
