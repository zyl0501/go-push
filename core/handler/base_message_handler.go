package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
)

type baseMessageHandler struct {
	wrap BaseMessageHandlerWrap
}

func (handler *baseMessageHandler) Handle(packet protocol.Packet, conn api.Conn) {
	msg := handler.wrap.Decode(packet, conn)
	if &msg == nil {
		msg.DecodeBody()
		handler.wrap.HandleMessage(msg)
	}
}

type BaseMessageHandlerWrap interface {
	Decode(packet protocol.Packet, connection api.Conn) api.Message

	HandleMessage(msg api.Message)
}
