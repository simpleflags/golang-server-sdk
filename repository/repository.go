package repository

import (
	"fmt"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/log"
)

// Callback provides events when repository data being modified
type Callback interface {
	OnFlagStored(identifier string)
	OnFlagDeleted(identifier string)
	OnVariableStored(identifier string)
	OnVariableDeleted(identifier string)
}

// Repository holds cache and optionally offline data
type Repository struct {
	cache    Cache
	storage  Storage
	callback Callback
}

type Option func(r *Repository)

func WithStorage(storage Storage) Option {
	return func(r *Repository) {
		r.storage = storage
	}
}

func WithCallback(callback Callback) Option {
	return func(r *Repository) {
		r.callback = callback
	}
}

// New repository with only cache capability
func New(cache Cache, options ...Option) Repository {
	r := Repository{
		cache: cache,
	}

	for _, option := range options {
		option(&r)
	}

	return r
}

func (r Repository) getConfigurationAndCache(identifier string, cacheable bool) (evaluation.Configuration, error) {
	flagKey := formatFlagKey(identifier)
	flag, ok := r.cache.Get(flagKey)
	if ok {
		return flag.(evaluation.Configuration), nil
	}

	if r.storage != nil {
		var flag evaluation.Configuration
		err := r.storage.Get(flagKey, &flag)
		if err == nil && cacheable {
			r.cache.Set(flagKey, flag)
			return flag, nil
		}
	}
	return evaluation.Configuration{}, fmt.Errorf("%w with identifier: %s", ErrFeatureConfigNotFound, identifier)
}

// GetConfiguration returns flag from cache or offline storage
func (r Repository) GetConfiguration(identifier string) (evaluation.Configuration, error) {
	return r.getConfigurationAndCache(identifier, true)
}

func (r Repository) getVariableAndCache(identifier string, cacheable bool) (evaluation.Variable, error) {
	variableKey := formatVariableKey(identifier)
	variable, ok := r.cache.Get(variableKey)
	if ok {
		return variable.(evaluation.Variable), nil
	}

	if r.storage != nil {
		var variable evaluation.Variable
		err := r.storage.Get(variableKey, &variable)
		if err == nil && cacheable {
			r.cache.Set(variableKey, variable)
			return variable, nil
		}
	}
	return evaluation.Variable{}, fmt.Errorf("%w with identifier: %s", ErrSegmentNotFound, identifier)
}

// GetVariable return variable from repository
func (r Repository) GetVariable(identifier string) (evaluation.Variable, error) {
	return r.getVariableAndCache(identifier, true)
}

// SetConfiguration places a flag in the repository with the new value
func (r Repository) SetConfiguration(config *evaluation.Configuration) {
	if r.isFlagOutdated(config) {
		return
	}
	flagKey := formatFlagKey(config.Identifier)
	if r.storage != nil {
		if err := r.storage.Set(flagKey, *config); err != nil {
			log.Errorf("error while storing the flag %s into repository", config.Identifier)
		}
		r.cache.Remove(flagKey)
	} else {
		r.cache.Set(flagKey, *config)
	}

	if r.callback != nil {
		r.callback.OnFlagStored(config.Identifier)
	}
}

// SetVariable places a variable in the repository with the new value
func (r Repository) SetVariable(variable *evaluation.Variable) {
	variableKey := formatVariableKey(variable.Identifier)
	if r.storage != nil {
		if err := r.storage.Set(variableKey, *variable); err != nil {
			log.Errorf("error while storing the variable %s into repository", variable.Identifier)
		}
		r.cache.Remove(variableKey)
	} else {
		r.cache.Set(variableKey, *variable)
	}

	if r.callback != nil {
		r.callback.OnVariableStored(variable.Identifier)
	}
}

// DeleteConfiguration removes a flag from the repository
func (r Repository) DeleteConfiguration(identifier string) {
	flagKey := formatFlagKey(identifier)
	if r.storage != nil {
		// remove from storage
		if err := r.storage.Remove(flagKey); err != nil {
			log.Errorf("error while removing flag %s from repository", identifier)
		}
	}
	// remove from cache
	r.cache.Remove(flagKey)
	if r.callback != nil {
		r.callback.OnFlagDeleted(identifier)
	}
}

// DeleteVariable removes a segment from the repository
func (r Repository) DeleteVariable(identifier string) {
	groupKey := formatVariableKey(identifier)
	if r.storage != nil {
		// remove from storage
		if err := r.storage.Remove(groupKey); err != nil {
			log.Errorf("error while removing target group %s from repository", identifier)
		}
	}
	// remove from cache
	r.cache.Remove(groupKey)
	if r.callback != nil {
		r.callback.OnVariableDeleted(identifier)
	}
}

func (r Repository) isFlagOutdated(config *evaluation.Configuration) bool {
	oldFlag, err := r.getConfigurationAndCache(config.Identifier, false)
	if err != nil {
		return false
	}

	return oldFlag.Version >= config.Version
}

// Close all resources
func (r Repository) Close() {

}

func formatFlagKey(identifier interface{}) string {
	return "flag__" + identifier.(string)
}

func formatVariableKey(identifier interface{}) string {
	return "variable__" + identifier.(string)
}
