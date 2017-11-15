package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
)

type PushUpHandler struct {
	*BaseMessageHandler
}

func NewPushUpHandler() *PushUpHandler{
	baseHandler := &BaseMessageHandler{}
	handler := PushUpHandler{BaseMessageHandler: baseHandler}
	handler.BaseMessageHandlerWrap = &handler
	return &handler
}

func (handler *PushUpHandler) Decode(packet protocol.Packet, conn api.Conn) api.Message {
	msg := message.NewPushUpMessage0(conn)
	msg.Pkt = packet
	return msg
}

func (handler *PushUpHandler) HandleMessage(m api.Message) {
	msg := m.(*message.PushUpMessage)
	log.Info("receive push request " + string(msg.Content))
}
