package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	log "github.com/alecthomas/log4go"
	"time"
)

type HeartBeatHandler struct {
}

func (handler *HeartBeatHandler) Handle(packet protocol.Packet, conn api.Conn) {
	conn.Send(packet); //ping -> pong
	conn.GetConn().SetDeadline(time.Now().Add(time.Duration(conn.GetSessionContext().Heartbeat)))
	log.Info("ping -> pong, %v", conn);
}
