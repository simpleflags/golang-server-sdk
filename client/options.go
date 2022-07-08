package client

import (
	"github.com/simpleflags/golang-server-sdk/repository"
)

// ConfigOption is used as return value for advanced client configuration
// using options pattern
type ConfigOption func(config *config)

// WithPullInterval set pulling interval in minutes
func WithPullInterval(interval uint) ConfigOption {
	return func(config *config) {
		config.pullInterval = interval
	}
}

// WithCache set custom cache or predefined one from cache package
func WithCache(cache repository.Cache) ConfigOption {
	return func(config *config) {
		config.cache = cache
	}
}

// WithStorage set custom storage
func WithStorage(storage repository.Storage) ConfigOption {
	return func(config *config) {
		config.storage = storage
	}
}

// WithPullerEnabled set puller on or off
func WithPullerEnabled(val bool) ConfigOption {
	return func(config *config) {
		config.enablePuller = val
	}
}

// WithStreamEnabled set stream on or off
func WithStreamEnabled(val bool) ConfigOption {
	return func(config *config) {
		config.enableStream = val
	}
}

// WithPrefetchFlags set of flags to be prefetched
func WithPrefetchFlags(identifiers ...string) ConfigOption {
	return func(config *config) {
		config.flags = identifiers
	}
}
