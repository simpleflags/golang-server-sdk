package client

import (
	"context"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/connector"
	"github.com/simpleflags/golang-server-sdk/log"
	"github.com/simpleflags/golang-server-sdk/repository"
	"go.uber.org/atomic"
	"time"
)

type puller struct {
	interval    uint
	connector   connector.Connector
	repository  repository.Repository
	stopped     chan struct{}
	init        *atomic.Bool
	identifiers []string
}

func newPuller(connector connector.Connector, repository repository.Repository, interval uint,
	identifiers ...string) puller {

	return puller{
		interval:    interval,
		connector:   connector,
		repository:  repository,
		stopped:     make(chan struct{}),
		init:        atomic.NewBool(false),
		identifiers: identifiers,
	}
}

func (p puller) flags(ctx context.Context) (evaluation.Configurations, error) {
	configurations, err := p.connector.Configurations(ctx, p.identifiers...)
	if err != nil {
		return evaluation.Configurations{}, err
	}

	for _, config := range configurations {
		p.repository.SetConfiguration(&config)
	}
	return configurations, nil
}

func (p puller) variables(ctx context.Context, identifiers ...string) {
	variables, err := p.connector.Variables(ctx, identifiers...)
	if err != nil {
		log.Errorf("error loading variables from server %v", err)
	}

	for _, variable := range variables {
		p.repository.SetVariable(&variable)
	}
}

func (p puller) initialized() bool {
	return p.init.Load()
}

func (p puller) pull(ctx context.Context) {
	log.Info("puller iteration")

	// first load flags from server
	configs, err := p.flags(ctx)
	if err != nil {
		log.Errorf("error loading flags from server %v", err)
	}

	// extract all variables from fetched flags
	variables := make([]string, 0)
	for _, cnf := range configs {
		for _, rule := range cnf.Rules {
			vars, err := evaluation.Variables(rule.Expression)
			if err != nil {
				return
			}
			variables = append(variables, vars...)
		}
	}
	// load variables
	p.variables(ctx, variables...)

	if !p.init.Load() {
		p.init.Store(true)
		log.Info("puller initialized")
	}
}

func (p puller) start(ctx context.Context) {
	log.Info("Starting puller")

	pullingTicker := time.NewTicker(time.Second * time.Duration(p.interval))
	go func() {
		p.pull(ctx)
		for {
			select {
			case <-p.stopped:
				return
			case <-ctx.Done():
				pullingTicker.Stop()
				return
			case <-pullingTicker.C:
				go p.pull(ctx)
			}
		}
	}()

	log.Info("Poller started")
}

func (p puller) stop() {
	log.Info("Stopping puller")
	p.stopped <- struct{}{}
	log.Info("Poller stopped")
}

func (p puller) close() {
	log.Info("Closing puller")
	p.stop()
	close(p.stopped)
	log.Info("Poller closed")
}
