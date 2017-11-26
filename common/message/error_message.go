package message

import (
	"io"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
)

type ErrorMessage struct {
	ByteBufMessage

	Cmd    byte
	Code   byte
	Reason string
	Data   string
}

func (message *ErrorMessage) DecodeByteBufMessage(reader io.Reader) {
	message.Cmd = DecodeByte(reader)
	message.Code = DecodeByte(reader)
	message.Reason = DecodeString(reader)
	message.Data = DecodeString(reader)
}

func (message *ErrorMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeByte(writer, message.Cmd)
	EncodeByte(writer, message.Code)
	EncodeString(writer, message.Reason)
	EncodeString(writer, message.Data)
}

func (msg *ErrorMessage) Send() {
	msg.sendRaw()
}

func NewErrorMessage(msg api.Message) *ErrorMessage {
	result := ErrorMessage{}
	packet := protocol.Packet{Cmd: protocol.ERROR, SessionId: msg.GetPacket().SessionId}
	conn := msg.GetConnection()
	result.Code = protocol.ERROR
	result.Cmd = packet.Cmd
	result.ByteBufMessage = ByteBufMessage{BaseMessage: &BaseMessage{Pkt: packet, Connection: conn}, ByteBufMessageCodec: &result}
	return &result
}

func NewErrorMessage0(sessionId uint32, conn api.Conn) *ErrorMessage {
	result := ErrorMessage{}
	packet := protocol.Packet{Cmd: protocol.ERROR, SessionId: sessionId}
	result.Code = protocol.ERROR
	result.Cmd = packet.Cmd
	result.ByteBufMessage = ByteBufMessage{BaseMessage: &BaseMessage{Pkt: packet, Connection: conn}, ByteBufMessageCodec: &result}
	return &result
}
