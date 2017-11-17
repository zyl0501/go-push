package handler

import (
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/common/message"
	"github.com/zyl0501/go-push/api/router"
)

type BindUserHandler struct {
	*BaseMessageHandler

	routerManager *router.LocalRouterManager
}

func NewBindUserHandler(routerManager *router.LocalRouterManager) *BindUserHandler {
	baseHandler := &BaseMessageHandler{}
	handler := BindUserHandler{BaseMessageHandler: baseHandler, routerManager: routerManager}
	handler.BaseMessageHandlerWrap = &handler
	return &handler
}

func (handler *BindUserHandler) Decode(packet protocol.Packet, conn api.Conn) api.Message {
	msg := message.NewBindUserMessage(packet, conn)
	return msg
}

func (handler *BindUserHandler) HandleMessage(m api.Message) {
	msg := m.(*message.BindUserMessage)
	log.Info("bind user")
	handler.routerManager.Register(msg.UserId, router.LocalRouter{Conn: msg.Connection})
}
