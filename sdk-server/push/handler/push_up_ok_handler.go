package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/core/handler"
)

type PushUpOKHandler struct {
	*handler.BaseMessageHandler
}

func (handler *PushUpOKHandler) Decode(packet protocol.Packet, conn api.Conn) api.Message {
	return message.NewOKMessage(packet, conn)
}

func (handler *PushUpOKHandler) HandleMessage(m api.Message) {
	msg := m.(*message.OKMessage)
	log.Info("push success " + msg.Data)
}
