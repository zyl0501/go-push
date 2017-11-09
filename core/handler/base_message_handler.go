package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
)


func Handle(handler BaseMessageHandler, packet protocol.Packet, conn api.Conn) {
	msg := handler.Decode(packet, conn)
	if &msg == nil {
		msg.DecodeBody()
		handler.HandleMessage(msg)
	}
}

type BaseMessageHandler interface {
	Decode(packet protocol.Packet, connection api.Conn) api.Message

	HandleMessage(msg api.Message)
}
