package core

import (
	api "github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/core/service"
	"github.com/zyl0501/go-push/core/push"
)

type MPushServer struct {
	ConnectionServer service.ConnectionServer
	PushCenter *push.PushCenter
	SessionManager *session.ReusableSessionManager
}

func NewPushServer() *MPushServer {
	pushServer := MPushServer{SessionManager:session.NewReusableSessionManager()}
	pushCenter := push.NewPushCenter()
	connectionServer := service.NewConnectionServer(pushServer.SessionManager,pushCenter)
	pushServer.ConnectionServer = connectionServer
	pushServer.PushCenter = pushCenter
	return &pushServer
}

func (pushServer *MPushServer) Init() {
	pushServer.ConnectionServer.Init()
}

func (pushServer *MPushServer) Start() {
	var listener api.Listener
	pushServer.PushCenter.Start()
	pushServer.ConnectionServer.Start(listener)
}

func (pushServer *MPushServer) Stop() {
	var listener api.Listener
	pushServer.ConnectionServer.Stop(listener)
}
