package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	log "github.com/alecthomas/log4go"
)

type HeartBeatHandler struct {
}

func (handler *HeartBeatHandler) Handle(packet protocol.Packet, conn api.Conn) {
	conn.Send(packet); //ping -> pong
	log.Info("ping -> pong, %v", conn);
}
