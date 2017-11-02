package message

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	"net"
)

type BaseMessage struct {
	message api.Message
	packet protocol.Packet
}

func (message *BaseMessage) getConnection() net.Conn {
	return nil
}

func (message *BaseMessage) decodeBody() {

}

func (message *BaseMessage) encodeBody() {

}

func (message *BaseMessage) getPacket() protocol.Packet {
	return message.packet
}
