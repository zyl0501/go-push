package api

import (
	"net"
	"github.com/zyl0501/go-push/api/protocol"
)

type Message interface {
	getConnection() net.Conn

	decodeBody()

	encodeBody()

	//send(listener ChannelFutureListener)

	getPacket() protocol.Packet
}

type MessageHandler interface {
	Handle(packet protocol.Packet, conn Conn) error
}

type PacketReceiver interface {
	OnReceive(packet protocol.Packet, conn Conn)
}