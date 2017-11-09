package message

import (
	"github.com/zyl0501/go-push/api"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
)

type BaseMessage struct {
	Pkt        protocol.Packet
	Connection api.Conn
	codec      baseMessageCodec
}

func (message *BaseMessage) GetConnection() api.Conn {
	return message.Connection
}

func (message *BaseMessage) DecodeBody() {
	message.decodeBaseMessage(message.codec, message)
}

func (message *BaseMessage) EncodeBody() {
	message.encodeBaseMessage(message.codec, message)
}

func (message *BaseMessage) GetPacket() protocol.Packet {
	return message.Pkt
}

func (message *BaseMessage) Send() {
	send(message)
}

func (message *BaseMessage) decodeBaseMessage(codec baseMessageCodec, m api.Message) {
	msg := m
	packet := msg.GetPacket()

	//1.解密
	tmp := packet.Body;
	//2.解压

	if len(tmp) == 0 {
		//"message decode ex"
		return
	}

	packet.Body = tmp
	codec.DecodeBaseMessage(packet.Body)
	packet.Body = nil // 释放内存
}

func (message *BaseMessage) encodeBaseMessage(codec baseMessageCodec, m api.Message) {
	tmp := codec.EncodeBaseMessage();
	if len(tmp) > 0 {
		//1.压缩
		//2.加密
	}
}

func send(msg api.Message) {
	msg.EncodeBody()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

func newBaseMessage(packet protocol.Packet, conn api.Conn, codec baseMessageCodec) *BaseMessage {
	return &BaseMessage{Pkt: packet, Connection: conn, codec: codec}
}

type baseMessageCodec interface {
	DecodeBaseMessage(body []byte)

	EncodeBaseMessage() ([]byte)
}
