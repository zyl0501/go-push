package core

import (
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/core/service"
	"github.com/zyl0501/go-push/core/push"
	"github.com/zyl0501/go-push/api/router"
)

type MPushServer struct {
	ConnectionServer service.ConnectionServer
	GatewayServer    *service.GatewayServer
	PushCenter       *push.PushCenter
}

func NewPushServer() *MPushServer {
	routerManager := router.NewLocalRouterManager()
	pushCenter := push.NewPushCenter(routerManager)
	pushServer := MPushServer{}
	pushServer.ConnectionServer = service.NewConnectionServer(session.NewReusableSessionManager(), routerManager)
	pushServer.GatewayServer = service.NewGatewayServer(pushCenter)
	pushServer.PushCenter = pushCenter
	return &pushServer
}

func (pushServer *MPushServer) Init() {
	pushServer.PushCenter.Init()
	pushServer.ConnectionServer.Init()
	pushServer.GatewayServer.Init()
}

func (pushServer *MPushServer) Start() {
	go pushServer.PushCenter.Start(nil)
	go pushServer.ConnectionServer.Start(nil)
	go pushServer.GatewayServer.Start(nil)
}

func (pushServer *MPushServer) Stop() {
	pushServer.GatewayServer.Stop(nil)
	pushServer.ConnectionServer.Stop(nil)
	pushServer.PushCenter.Stop(nil)
}
