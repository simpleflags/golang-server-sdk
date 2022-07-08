package sfsdk

import (
	"errors"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/client"
	"github.com/simpleflags/golang-server-sdk/connector"
	"github.com/simpleflags/golang-server-sdk/log"
	"sync"
)

var (
	once          sync.Once
	defaultClient client.Client
)

func Initialize(apiKey string, options ...client.ConfigOption) error {
	var err error
	once.Do(func() {
		defaultClient, err = client.New(apiKey, options...)
	})
	return err
}

func InitWithConnector(connector connector.Connector, options ...client.ConfigOption) error {
	var err error
	once.Do(func() {
		defaultClient, err = client.NewWithConnector(connector, options...)
	})
	return err
}

func WaitForInitialization() {
	if defaultClient != nil {
		defaultClient.WaitForInitialization()
	}
}

func Evaluate(feature string, target evaluation.Target) evaluation.Evaluation {
	if defaultClient != nil {
		return defaultClient.Evaluate(feature, target)
	}
	return evaluation.Evaluation{}
}

func Close() error {
	if defaultClient != nil {
		return defaultClient.Close()
	}
	return errors.New("defaultClient was not initialized")
}

// SetLogger sets the default logger to be used by this package
func SetLogger(logger log.Logger) {
	log.SetLogger(logger)
}
