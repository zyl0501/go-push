package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"io"
)

type HandshakeMessage struct {
	Pkt        protocol.Packet
	Connection api.Conn

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

func (message *HandshakeMessage) GetConnection() api.Conn {
	return message.Connection
}

func (message *HandshakeMessage) GetPacket() protocol.Packet {
	return message.Pkt
}

func (message *HandshakeMessage) Send() {
	Send(message)
}

func (message *HandshakeMessage) DecodeByteBufMessage(reader io.Reader) {
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

func (message *HandshakeMessage) EncodeByteBufMessage(writer io.Writer) {
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

func (message *HandshakeMessage) DecodeBody() {
	DecodeBaseMessage(message, message)
}

func (message *HandshakeMessage) EncodeBody() {
	EncodeBaseMessage(message, message)
}

func (message *HandshakeMessage) DecodeBaseMessage(body []byte) {
	DecodeByteBufMessage(message, body)
}

func (message *HandshakeMessage) EncodeBaseMessage() ([]byte) {
	return EncodeByteBufMessage(message, message)
}
