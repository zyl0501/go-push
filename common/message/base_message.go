package message

import (
	"github.com/zyl0501/go-push/api"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
)

type BaseMessage struct {
	baseMessageCodec
	Pkt        protocol.Packet
	Connection api.Conn
}

func (msg *BaseMessage) GetConnection() api.Conn {
	return msg.Connection
}

func (msg *BaseMessage) DecodeBody() {
	packet := msg.GetPacket()


	tmp := packet.Body;
	//1.解密
	//2.解压

	if len(tmp) == 0 {
		//"message decode ex"
		return
	}

	packet.Body = tmp
	msg.decodeBaseMessage(packet.Body)
	packet.Body = nil // 释放内存
}

func (msg *BaseMessage) EncodeBody() {
	tmp := msg.encodeBaseMessage();
	if len(tmp) > 0 {
		//1.压缩
		//2.加密

		msg.Pkt.Body = tmp
	}
}

func (msg *BaseMessage) GetPacket() protocol.Packet {
	return msg.Pkt
}

func (msg *BaseMessage) Send() {
	msg.EncodeBody()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

func (msg *ByteBufMessage) sendRaw() {
	msg.encodeRaw()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

func (msg *ByteBufMessage) encodeRaw() {
	tmp := msg.encodeBaseMessage()
	msg.Pkt.Body = tmp
}

type baseMessageCodec interface {
	decodeBaseMessage(body []byte)
	encodeBaseMessage() ([]byte)
}
