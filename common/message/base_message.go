package message

import (
	"github.com/zyl0501/go-push/api"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
	log "github.com/alecthomas/log4go"
	"errors"
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
	if packet.HasFlag(protocol.FLAG_CRYPTO) {
		cip := msg.Connection.GetSessionContext().Cipher0
		if cip != nil {
			var err error
			tmp, err = cip.Decrypt(tmp);
			if err != nil {
				log.Error("decrypt error ", err)
				return
			}
		}
	}
	//2.解压

	if len(tmp) == 0 {
		log.Error(errors.New("message decode ex"))
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
		context := msg.Connection.GetSessionContext()
		if context.Cipher0 != nil {
			result, _ := context.Cipher0.Encrypt(tmp);
			if len(result) > 0 {
				tmp = result;
				msg.Pkt.AddFlag(protocol.FLAG_CRYPTO);
			}
		}
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
	writer.Flush()
}

func (msg *ByteBufMessage) sendRaw() {
	msg.encodeRaw()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
	writer.Flush()
}

func (msg *ByteBufMessage) encodeRaw() {
	tmp := msg.encodeBaseMessage()
	msg.Pkt.Body = tmp
}

type baseMessageCodec interface {
	decodeBaseMessage(body []byte)
	encodeBaseMessage() ([]byte)
}
