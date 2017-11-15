package core

import (
	api "github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/core/service"
)

type MPushServer struct {
	ConnectionServer service.ConnectionServer
	SessionManager *session.ReusableSessionManager
}

func NewPushServer() *MPushServer {
	pushServer := MPushServer{SessionManager:session.NewReusableSessionManager()}
	connectionServer := service.NewConnectionServer(pushServer.SessionManager)
	pushServer.ConnectionServer = connectionServer
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
