package message

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	"net"
	"github.com/zyl0501/go-push/core/service"
	"github.com/zyl0501/go-push/core/connection"
)

type BaseMessage struct {
	packet protocol.Packet
	connection api.Conn
}

func (message *BaseMessage) GetConnection() api.Conn {
	return message.connection
}

func (message *BaseMessage) DecodeBody() {

}

func (message *BaseMessage) EncodeBody() {

}

func (message *BaseMessage) GetPacket() protocol.Packet {
	return message.packet
}
