package client

import (
	"github.com/simpleflags/evaluation"
)

type Client interface {
	WaitForInitialization()
	Evaluate(feature string, target evaluation.Target) evaluation.Evaluation
	Close() error
}
