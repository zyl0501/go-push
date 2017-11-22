package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"io"
	"time"
	"github.com/zyl0501/go-push/tools/utils"
)

type HandshakeMessage struct {
	*ByteBufMessage

	DeviceId      string
	OsName        string
	OsVersion     string
	ClientVersion string
	Iv            []byte
	ClientKey     []byte
	MinHeartbeat  time.Duration
	MaxHeartbeat  time.Duration
	Timestamp     int64
}

func NewHandshakeMessage(packet protocol.Packet, conn api.Conn) *HandshakeMessage {
	//pkt := protocol.Packet{Cmd:protocol.HANDSHAKE, SessionId:packet.SessionId}
	baseMessage := BaseMessage{Pkt:packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := HandshakeMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func NewHandshakeMessage0(conn api.Conn) *HandshakeMessage {
	packet := protocol.Packet{Cmd:protocol.HANDSHAKE, SessionId:protocol.GetSessionId()}
	baseMessage := BaseMessage{Pkt:packet, Connection: conn}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage}
	msg := HandshakeMessage{ByteBufMessage: &byteMessage}
	msg.BaseMessageCodec = &msg
	msg.ByteBufMessageCodec = &msg
	return &msg
}

func (message *HandshakeMessage) DecodeByteBufMessage(reader io.Reader) {
	message.DeviceId = DecodeString(reader)
	message.OsName = DecodeString(reader)
	message.OsVersion = DecodeString(reader)
	message.ClientVersion = DecodeString(reader)
	message.Iv = DecodeBytes(reader)
	message.ClientKey = DecodeBytes(reader)
	message.MinHeartbeat =  utils.MillisecondToDuration(DecodeInt64(reader))
	message.MaxHeartbeat =  utils.MillisecondToDuration(DecodeInt64(reader))
	message.Timestamp = DecodeInt64(reader)
}

func (message *HandshakeMessage) EncodeByteBufMessage(writer io.Writer) {
	EncodeString(writer, message.DeviceId)
	EncodeString(writer, message.OsName)
	EncodeString(writer, message.OsVersion)
	EncodeString(writer, message.ClientVersion)
	EncodeBytes(writer, message.Iv)
	EncodeBytes(writer, message.ClientKey)
	EncodeInt64(writer, utils.DurationToMillisecond(message.MinHeartbeat))
	EncodeInt64(writer, utils.DurationToMillisecond(message.MaxHeartbeat))
	EncodeInt64(writer, message.Timestamp)
}
