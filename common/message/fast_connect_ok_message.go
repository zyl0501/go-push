package message

import (
	"io"
	"github.com/zyl0501/go-push/tools/utils"
	"time"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
)

type FastConnectOKMessage struct {
	*ByteBufMessage

	Heartbeat time.Duration
}

func NewFastConnectOKMessage(sessionId uint32, conn api.Conn) *FastConnectOKMessage {
	pkt := protocol.Packet{Cmd:protocol.FAST_CONNECT, SessionId:sessionId}
	baseMessage := BaseMessage{Pkt: pkt, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := FastConnectOKMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func NewFastConnectOKMessage0(conn api.Conn) *FastConnectOKMessage {
	packet := protocol.Packet{Cmd: protocol.FAST_CONNECT, SessionId: protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt: packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := FastConnectOKMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *FastConnectOKMessage) DecodeByteBufMessage(reader io.Reader) {
	message.Heartbeat = utils.MillisecondToDuration(DecodeInt64(reader))
}

func (message *FastConnectOKMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeInt64(writer, utils.DurationToMillisecond(message.Heartbeat))
}

func (msg *FastConnectOKMessage) Send() {
	msg.sendRaw()
}
