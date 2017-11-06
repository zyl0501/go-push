package handler

import (
	"github.com/zyl0501/go-push/api"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/protocol"
)

type HandshakeHandler struct {
}

func (handler HandshakeHandler) Handle(packet protocol.Packet, conn api.Conn) error {
	log.Debug("HandshakeHandler invoke")
	return nil
}
