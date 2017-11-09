package api

import (
	"net"
	"github.com/zyl0501/go-push/api/protocol"
)

type Message interface {
	GetConnection() Conn

	DecodeBody()

	EncodeBody()

	//send(listener ChannelFutureListener)

	GetPacket() protocol.Packet
}

type MessageHandler interface {
	Handle(packet protocol.Packet, conn Conn)
}

type PacketReceiver interface {
	OnReceive(packet protocol.Packet, conn Conn)
}