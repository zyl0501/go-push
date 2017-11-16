package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/core/push"
)

type PushUpHandler struct {
	*BaseMessageHandler
	pushCenter *push.PushCenter
}

func NewPushUpHandler(pushCenter *push.PushCenter) *PushUpHandler{
	baseHandler := &BaseMessageHandler{}
	handler := PushUpHandler{BaseMessageHandler: baseHandler,pushCenter:pushCenter}
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

	handler.pushCenter.Push(*msg)
}
