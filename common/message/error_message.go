package message

import (
	"io"
	"github.com/zyl0501/go-push/api"
)

type ErrorMessage struct {
	ByteBufMessage

	Cmd    byte
	Code   byte
	Reason string
	Data   string
}

func (message *ErrorMessage) decodeByteBufMessage(reader io.Reader) {
	message.Cmd = DecodeByte(reader)
	message.Code = DecodeByte(reader)
	message.Reason = DecodeString(reader)
	message.Data = DecodeString(reader)
}

func (message *ErrorMessage) encodeByteBufMessage(writer io.Writer) {
	EncodeByte(writer, message.Cmd)
	EncodeByte(writer, message.Code)
	EncodeString(writer, message.Reason)
	EncodeString(writer, message.Data)
}

func NewErrorMessage(msg api.Message) *ErrorMessage {
	result := ErrorMessage{}

	packet := msg.GetPacket()
	conn := msg.GetConnection()
	result.Cmd = packet.Cmd
	result.ByteBufMessage = ByteBufMessage{Pkt: packet, Connection: conn, byteBufMessageCodec: &result}
	return &result
}
