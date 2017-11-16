package message

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
)

type PushDownMessage struct {
	*BaseMessage

	Content []byte
}

func NewPushDownMessage0(conn api.Conn) *PushDownMessage {
	packet := protocol.Packet{Cmd:protocol.PUSH, SessionId:protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt:packet, Connection: conn}
	msg := PushDownMessage{BaseMessage: &baseMessage}
	msg.BaseMessageCodec = &msg
	return &msg
}

func (msg *PushDownMessage) decodeBaseMessage(body []byte) {
	msg.Content = body
}

func (msg *PushDownMessage) encodeBaseMessage() ([]byte) {
	return msg.Content
}