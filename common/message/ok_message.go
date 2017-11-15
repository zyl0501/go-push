package message

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	"io"
)

type OKMessage struct {
	*ByteBufMessage

	Code byte
	Data string
}

func NewOKMessage(packet protocol.Packet, conn api.Conn) *OKMessage {
	Pkt := protocol.Packet{Cmd:protocol.OK, SessionId:packet.SessionId}
	baseMessage := BaseMessage{Pkt:Pkt, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := OKMessage{ByteBufMessage: &byteMessage}
	msg.baseMessageCodec = &msg
	msg.byteBufMessageCodec = &msg
	return &msg
}

func NewOKMessage0(conn api.Conn) *OKMessage {
	packet := protocol.Packet{Cmd:protocol.OK, SessionId:protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt:packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := OKMessage{ByteBufMessage: &byteMessage}
	msg.baseMessageCodec = &msg
	msg.byteBufMessageCodec = &msg
	return &msg
}

func (message *OKMessage) decodeByteBufMessage(reader io.Reader) {
	message.Code = DecodeByte(reader)
	message.Data = DecodeString(reader)
}

func (message *OKMessage) encodeByteBufMessage(writer io.Writer) {
	EncodeByte(writer, message.Code)
	EncodeString(writer, message.Data)
}