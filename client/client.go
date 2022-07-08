package client

import (
	"context"
	"errors"
	"github.com/looplab/fsm"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/connector"
	"github.com/simpleflags/golang-server-sdk/connector/simple"
	"github.com/simpleflags/golang-server-sdk/repository"
	"go.uber.org/atomic"
)

// client is the Feature Flag client.
//
// This object evaluates feature flags and communicates with Feature Flag services.
// Applications should instantiate a single instance for the lifetime of their application
// and share it wherever features flags need to be evaluated.
//
// When an application is shutting down or no longer needs to use the client instance, it
// should call Close() to ensure that all of its connections and goroutines are shut down and
// that any pending analytics events have been delivered.
//
type client struct {
	config     config
	evaluator  *evaluation.Evaluator
	repository repository.Repository
	connector  connector.Connector
	puller     *puller
	updater    *updater
	stop       chan struct{}
	stopped    *atomic.Bool
	state      *fsm.FSM
}

// New creates a new client instance that connects to CF with the default configuration.
// For advanced configuration options use ConfigOptions functions
func New(apiKey string, options ...ConfigOption) (*client, error) {

	if apiKey == "" {
		return &client{}, ErrSdkCantBeEmpty
	}

	sfConn := simple.NewHttpConnector(apiKey)

	return NewWithConnector(sfConn, options...)
}

func NewWithConnector(connector connector.Connector, options ...ConfigOption) (*client, error) {
	if connector == nil {
		return &client{}, ErrConnectorCannotBeNil
	}

	//  functional options for config
	config, err := newDefaultConfig()
	for _, opt := range options {
		opt(&config)
	}

	repo := repository.New(config.cache)
	if config.storage != nil {
		repo = repository.New(config.cache, repository.WithStorage(config.storage))
	}

	evaluator, err := evaluation.NewEvaluator(repo)
	if err != nil {
		return nil, err
	}

	p := newPuller(connector, repo, config.pullInterval, config.flags...)

	state := fsm.NewFSM("disconnected",
		fsm.Events{
			{Name: "connect", Src: []string{"disconnected"}, Dst: "connected"},
			{Name: "disconnect", Src: []string{"connected"}, Dst: "disconnected"},
		},
		fsm.Callbacks{
			"enter_connected": func(e *fsm.Event) {
				p.stop()
			},
			"enter_disconnected": func(e *fsm.Event) {
				//p.start()
			},
		},
	)
	u := newUpdater(connector, repo, state)

	client := &client{
		config:     config,
		repository: repo,
		connector:  connector,
		evaluator:  evaluator,
		puller:     &p,
		updater:    &u,
		stop:       make(chan struct{}),
		stopped:    atomic.NewBool(false),
		state:      state,
	}

	client.start()

	return client, nil
}

func (c *client) WaitForInitialization() {
	for {
		if c.puller.initialized() {
			break
		}
	}
}

func (c *client) start() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c.stop
		cancel()
	}()

	if c.config.enablePuller {
		c.puller.start(ctx)
	} else {
		// good for lambda and short living environments
		c.puller.pull(ctx)
	}

	if c.config.enableStream {
		c.updater.start(ctx)
	}
}

func (c *client) Evaluate(feature string, target evaluation.Target) evaluation.Evaluation {
	eval := c.evaluator.Evaluate(feature, target)
	//c.analyticsService.PushToQueue(feature, target, variation)
	return eval
}

// Close shuts down the Feature Flag client. After calling this, the client
// should no longer be used
func (c *client) Close() error {
	if c.stopped.Load() {
		return errors.New("client already closed")
	}
	close(c.stop)

	c.stopped.Store(true)
	return nil
}
