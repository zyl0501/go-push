package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"io"
)

type BindUserMessage struct {
	*ByteBufMessage

	UserId string
	Tags   string
	Data   string
}

func NewBindUserMessage(packet protocol.Packet, conn api.Conn) *BindUserMessage {
	baseMessage := BaseMessage{Pkt: packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := BindUserMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *BindUserMessage) DecodeByteBufMessage(reader io.Reader) {
	message.UserId = DecodeString(reader)
	message.Tags = DecodeString(reader)
	message.Data = DecodeString(reader)
}

func (message *BindUserMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeString(writer, message.UserId)
	EncodeString(writer, message.Tags)
	EncodeString(writer, message.Data)
}
