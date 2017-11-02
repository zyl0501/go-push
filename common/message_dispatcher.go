package common

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	log "github.com/alecthomas/log4go"
	"net"
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

func (dispatcher *MessageDispatcher) onReceive(packet protocol.Packet, conn net.Conn) {
	handler := dispatcher.handlers[packet.Cmd]
	if handler != nil {
		err := handler.Handle(packet, conn)
		if err != nil {
			log.Error("dispatch message ex, packet={}, connect={}, body={}")
		}
	} else {
		if dispatcher.unsupportedPolicy > POLICY_IGNORE {
			log.Error("dispatch message failure, cmd={} unsupported, packet={}, connect={}, body={}")
			if dispatcher.unsupportedPolicy == POLICY_REJECT {

			}
		}
	}
}
