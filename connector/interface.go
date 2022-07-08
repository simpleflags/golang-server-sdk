package connector

import (
	"context"
	"github.com/simpleflags/evaluation"
)

type Msg struct {
	Event []byte
	Data  []byte
}

type Updater interface {
	OnDisconnect()
	OnEvent(msg *Msg)
	OnConnect()
}

type Connector interface {
	Configurations(ctx context.Context, identifiers ...string) (evaluation.Configurations, error)
	Variables(ctx context.Context, identifiers ...string) ([]evaluation.Variable, error)
	Stream(ctx context.Context, updater Updater) error
	Close() error
}
