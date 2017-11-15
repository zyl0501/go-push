package common

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	log "github.com/alecthomas/log4go"
)

const (
	POLICY_IGNORE int = 0
	POLICY_LOG    int = 1
	POLICY_REJECT int = 2
)

type MessageDispatcher struct {
	handlers          map[byte]api.MessageHandler
	unsupportedPolicy int
}

func NewMessageDispatcher() (dispatcher MessageDispatcher) {
	return MessageDispatcher{handlers: make(map[byte]api.MessageHandler), unsupportedPolicy: POLICY_LOG}
}

func (dispatcher *MessageDispatcher) Register(cmd byte, handler api.MessageHandler) {
	dispatcher.handlers[cmd] = handler
}

func (dispatcher *MessageDispatcher) OnReceive(packet protocol.Packet, conn api.Conn) {
	handler := dispatcher.handlers[packet.Cmd]
	if handler != nil {
		handler.Handle(packet, conn)
	} else {
		if dispatcher.unsupportedPolicy > POLICY_IGNORE {
			log.Error("dispatch message failure, cmd={} unsupported, packet={}, connect={}, body={}")
			if dispatcher.unsupportedPolicy == POLICY_REJECT {
			}
		}
	}
}
