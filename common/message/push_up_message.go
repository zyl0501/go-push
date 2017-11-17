package message

import (
	"io"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
)

type PushUpMessage struct {
	*ByteBufMessage

	UserId string
	ClientType byte
	Timeout int64
	Content []byte
	Tags []string
}

func NewPushUpMessage(packet protocol.Packet, conn api.Conn) *PushUpMessage {
	Pkt := protocol.Packet{Cmd:protocol.PUSH_UP, SessionId:packet.SessionId}
	baseMessage := BaseMessage{Pkt:Pkt, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := PushUpMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func NewPushUpMessage0(conn api.Conn) *PushUpMessage {
	packet := protocol.Packet{Cmd:protocol.PUSH_UP, SessionId:protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt:packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := PushUpMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *PushUpMessage) DecodeByteBufMessage(reader io.Reader) {
	message.UserId = DecodeString(reader)
	message.ClientType = DecodeByte(reader)
	message.Timeout = DecodeInt64(reader)
	message.Content = DecodeBytes(reader)
	//message.Tags = DecodeBytes(reader)
}

func (message *PushUpMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeString(writer, message.UserId)
	EncodeByte(writer, message.ClientType)
	EncodeInt64(writer, message.Timeout)
	EncodeBytes(writer, message.Content)
	//EncodeBytes(writer, message.Tags)
}


func (msg *PushUpMessage) Send() {
	msg.sendRaw()
}

