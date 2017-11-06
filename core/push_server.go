package core

import (
	"github.com/zyl0501/go-push/core/service"
	api "github.com/zyl0501/go-push/api/service"
)

type MPushServer struct {
	ConnectionServer service.ConnectionServer
}

func NewPushServer() *MPushServer {
	connectionServer := service.NewConnectionServer()
	pushServer := MPushServer{connectionServer}
	return &pushServer
}

func (pushServer *MPushServer) Init() {
	pushServer.ConnectionServer.Init()
}

func (pushServer *MPushServer) Start() {
	var listener api.Listener
	pushServer.ConnectionServer.Start(listener)
}

func (pushServer *MPushServer) Stop() {
	var listener api.Listener
	pushServer.ConnectionServer.Stop(listener)
}
