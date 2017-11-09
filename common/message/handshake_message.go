package message

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"bytes"
	"bufio"
)

type HandshakeMessage struct {
	Pkt protocol.Packet
	Connection api.Conn

	byteBufMessage ByteBufMessage
	DeviceId       string
	OsName         string
	OsVersion      string
	ClientVersion  string
	Iv             []byte
	ClientKey      []byte
	MinHeartbeat   int32
	MaxHeartbeat   int32
	Timestamp      int64
}

func (message *HandshakeMessage) GetConnection() api.Conn {
	return message.byteBufMessage.GetConnection()
}

func (message *HandshakeMessage) DecodeBody() {
	message.byteBufMessage.DecodeBody()

	reader := bytes.NewReader(message.GetPacket().Body)
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

func (message *HandshakeMessage) EncodeBody() {
	message.byteBufMessage.EncodeBody()

	writer := bufio.NewWriter(message.GetConnection().GetConn())
	EncodeString(writer, message.DeviceId)
	EncodeString(writer, message.OsName)
	EncodeString(writer, message.OsVersion)
	EncodeString(writer, message.ClientVersion)
	EncodeBytes(writer, message.Iv)
	EncodeBytes(writer, message.ClientKey)
	EncodeInt32(writer, message.MinHeartbeat)
	EncodeInt32(writer, message.MaxHeartbeat)
	EncodeInt64(writer, message.Timestamp)

	writer.Flush()
}

func (message *HandshakeMessage) GetPacket() protocol.Packet {
	return message.byteBufMessage.GetPacket()
}
