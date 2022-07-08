package client

import (
	"context"
	"encoding/json"
	"github.com/looplab/fsm"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/connector"
	"github.com/simpleflags/golang-server-sdk/repository"
	"log"
)

type updater struct {
	connector  connector.Connector
	repository repository.Repository
	msgChannel chan *connector.Msg
	fsm        *fsm.FSM
}

func newUpdater(conn connector.Connector, repo repository.Repository, fsm *fsm.FSM) updater {
	msgChannel := make(chan *connector.Msg, 100)
	return updater{
		connector:  conn,
		repository: repo,
		msgChannel: msgChannel,
		fsm:        fsm,
	}
}

func (u *updater) start(ctx context.Context) {
	err := u.connector.Stream(ctx, u)
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		go u.consumer(ctx, u.msgChannel)
	}
}

func (u *updater) OnConnect() {
	err := u.fsm.Event("connect")
	if err != nil {
		log.Println(err)
	}
}

func (u *updater) OnDisconnect() {
	err := u.fsm.Event("disconnect")
	if err != nil {
		log.Println(err)
	}
}

func (u *updater) OnEvent(msg *connector.Msg) {
	u.msgChannel <- msg
}

func (u *updater) consumer(ctx context.Context, msgChan chan *connector.Msg) {
	for msg := range msgChan {
		ev := string(msg.Event)
		switch ev {
		case evaluation.CreateFlagEvent, evaluation.PatchFlagEvent:
			var cfg evaluation.Configuration
			err := json.Unmarshal(msg.Data, &cfg)
			if err != nil {
				log.Printf("error processing event %v", err)
				return
			}
			u.repository.SetConfiguration(&cfg)
		case evaluation.DeleteFlagEvent:
			u.repository.DeleteConfiguration(string(msg.Data))
		case evaluation.CreateVariable, evaluation.PatchVariable:
			var v evaluation.Variable
			err := json.Unmarshal(msg.Data, &v)
			if err != nil {
				log.Printf("error processing event %v", err)
				return
			}
			u.repository.SetVariable(&v)
		case evaluation.DeleteVariable:
			u.repository.DeleteVariable(string(msg.Data))
		}
	}
	<-ctx.Done()
}

func (u *updater) close() {
	close(u.msgChannel)
}
