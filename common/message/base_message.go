package message

import (
	"github.com/zyl0501/go-push/api"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
)

func DecodeBaseMessage(codec BaseMessageCodec, m api.Message) {
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

func EncodeBaseMessage(codec BaseMessageCodec, m api.Message) {
	tmp := codec.EncodeBaseMessage();
	if len(tmp) > 0 {
		//1.压缩
		//2.加密
	}
}

func Send(msg api.Message) {
	msg.EncodeBody()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

type BaseMessageCodec interface {
	DecodeBaseMessage(body []byte)

	EncodeBaseMessage() ([]byte)
}
