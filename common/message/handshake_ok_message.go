package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"io"
	"time"
	"github.com/zyl0501/go-push/tools/utils"
)

type HandshakeOKMessage struct {
	*ByteBufMessage

	ServerKey  []byte
	Heartbeat  time.Duration
	SessionId  string
	ExpireTime int64
}

func NewHandshakeOKMessage(packet protocol.Packet, conn api.Conn) *HandshakeOKMessage {
	pkt := protocol.Packet{Cmd: protocol.HANDSHAKE, SessionId: packet.SessionId}
	baseMessage := BaseMessage{Pkt: pkt, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := HandshakeOKMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func NewHandshakeOKMessage0(conn api.Conn) *HandshakeOKMessage {
	packet := protocol.Packet{Cmd: protocol.HANDSHAKE, SessionId: protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt: packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := HandshakeOKMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *HandshakeOKMessage) DecodeByteBufMessage(reader io.Reader) {
	message.ServerKey = DecodeBytes(reader)
	message.Heartbeat = utils.MillisecondToDuration(DecodeInt64(reader))
	message.SessionId = DecodeString(reader)
	message.ExpireTime = DecodeInt64(reader)
}

func (message *HandshakeOKMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeBytes(writer, message.ServerKey)
	EncodeInt64(writer, utils.DurationToMillisecond(message.Heartbeat))
	EncodeString(writer, message.SessionId)
	EncodeInt64(writer, message.ExpireTime)
}
