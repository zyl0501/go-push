package message

import (
	"testing"
	"github.com/zyl0501/go-push/api/protocol"
)

func TestNewHandshakeMessage(t *testing.T) {
	packet := protocol.Packet{}
	packet.Cmd = protocol.PUSH
	packet.Cc = 22
	packet.Lrc = 4
	packet.Flags = 2
	packet.SessionId = 234
	msg := NewHandshakeMessage(packet, nil)
	msg.OsVersion = "1.0.1"
	msg.OsName = "android"
	msg.EncodeBody()

	msg.OsVersion = ""
	msg.OsName = ""
	msg.DecodeBody()
	if msg.OsName == "android" && msg.OsVersion == "1.0.1"{
		t.Log("OK")
	}else{
		t.Error("Failure")
	}
}
