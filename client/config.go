package client

import (
	"github.com/simpleflags/golang-server-sdk/repository"
)

type config struct {
	pullInterval    uint // in seconds
	pushInterval    uint // in seconds
	cache           repository.Cache
	storage         repository.Storage
	enablePuller    bool
	enableStream    bool
	enableAnalytics bool
	flags           []string
}

func newDefaultConfig() (config, error) {

	defaultCache, err := repository.NewLruCache(1000)
	if err != nil {
		return config{}, err
	}

	return config{
		pullInterval:    60,
		pushInterval:    60,
		cache:           defaultCache,
		enablePuller:    true,
		enableStream:    true,
		enableAnalytics: true,
		flags:           []string{},
	}, nil
}
