package client

import "errors"

var (
	ErrSdkCantBeEmpty       = errors.New("SDK key cannot be empty")
	ErrConnectorCannotBeNil = errors.New("connector cannot be nil")
)
