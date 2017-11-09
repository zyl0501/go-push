package message

import (
	"testing"
	"github.com/zyl0501/go-push/api/protocol"
)

func TestNewHandshakeMessage(t *testing.T) {
	packet := protocol.Packet{}
	msg := NewHandshakeMessage(packet,nil)
	msg.ByteBufMessage.DecodeBody()
	msg.DecodeByteBufMessage(nil)

	msg2 := ByteBufMessage{}
	msg2.DecodeBody()
}
