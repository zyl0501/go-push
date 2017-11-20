package message

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/api/protocol"
	"io"
)

type OKMessage struct {
	*ByteBufMessage

	Cmd byte
	Code byte
	Data string
}

func NewOKMessage(packet protocol.Packet, conn api.Conn) *OKMessage {
	Pkt := protocol.Packet{Cmd:protocol.OK, SessionId:packet.SessionId}
	baseMessage := BaseMessage{Pkt:Pkt, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := OKMessage{ByteBufMessage: &byteMessage}
	msg.Cmd = packet.Cmd
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *OKMessage) DecodeByteBufMessage(reader io.Reader) {
	message.Cmd = DecodeByte(reader)
	message.Code = DecodeByte(reader)
	message.Data = DecodeString(reader)
}

func (message *OKMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeByte(writer, message.Cmd)
	EncodeByte(writer, message.Code)
	EncodeString(writer, message.Data)
}

//func (msg *OKMessage) Send() {
//	msg.sendRaw()
//}