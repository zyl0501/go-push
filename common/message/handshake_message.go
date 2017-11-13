package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"io"
)

type HandshakeMessage struct {
	*ByteBufMessage

	DeviceId      string
	OsName        string
	OsVersion     string
	ClientVersion string
	Iv            []byte
	ClientKey     []byte
	MinHeartbeat  int32
	MaxHeartbeat  int32
	Timestamp     int64
}

func NewHandshakeMessage(packet protocol.Packet, conn api.Conn) *HandshakeMessage {
	baseMessage := BaseMessage{}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage, Pkt: packet, Connection: conn}
	msg := HandshakeMessage{ByteBufMessage: &byteMessage}
	msg.baseMessageCodec = &msg
	msg.byteBufMessageCodec = &msg
	return &msg
}

func NewHandshakeMessage(conn api.Conn) *HandshakeMessage {
	baseMessage := BaseMessage{}
	byteMessage := ByteBufMessage{BaseMessage: &baseMessage, Pkt: packet, Connection: conn}
	msg := HandshakeMessage{ByteBufMessage: &byteMessage}
	msg.baseMessageCodec = &msg
	msg.byteBufMessageCodec = &msg
	return &msg
}

func (message *HandshakeMessage) decodeByteBufMessage(reader io.Reader) {
	message.DeviceId = DecodeString(reader)
	message.OsName = DecodeString(reader)
	message.OsVersion = DecodeString(reader)
	message.ClientVersion = DecodeString(reader)
	message.Iv = DecodeBytes(reader)
	message.ClientKey = DecodeBytes(reader)
	message.MinHeartbeat = DecodeInt32(reader)
	message.MaxHeartbeat = DecodeInt32(reader)
	message.Timestamp = DecodeInt64(reader)
}

func (message *HandshakeMessage) encodeByteBufMessage(writer io.Writer) {
	EncodeString(writer, message.DeviceId)
	EncodeString(writer, message.OsName)
	EncodeString(writer, message.OsVersion)
	EncodeString(writer, message.ClientVersion)
	EncodeBytes(writer, message.Iv)
	EncodeBytes(writer, message.ClientKey)
	EncodeInt32(writer, message.MinHeartbeat)
	EncodeInt32(writer, message.MaxHeartbeat)
	EncodeInt64(writer, message.Timestamp)
}
