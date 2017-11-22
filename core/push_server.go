package core

import (
	api "github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/core/service"
	"github.com/zyl0501/go-push/core/push"
	"github.com/zyl0501/go-push/api/router"
)

type MPushServer struct {
	ConnectionServer service.ConnectionServer
	GatewayServer    *service.GatewayServer
	PushCenter       *push.PushCenter
	SessionManager   *session.ReusableSessionManager
}

func NewPushServer() *MPushServer {
	routerManager := router.NewLocalRouterManager()
	pushServer := MPushServer{SessionManager: session.NewReusableSessionManager()}
	pushCenter := push.NewPushCenter(routerManager)
	pushServer.ConnectionServer = service.NewConnectionServer(pushServer.SessionManager, pushCenter, routerManager)
	pushServer.GatewayServer = service.NewGatewayServer(pushCenter)
	pushServer.PushCenter = pushCenter
	return &pushServer
}

func (pushServer *MPushServer) Init() {
	pushServer.ConnectionServer.Init()
	pushServer.GatewayServer.Init()
}

func (pushServer *MPushServer) Start() {
	var listener api.Listener
	go pushServer.PushCenter.Start()
	go pushServer.ConnectionServer.Start(listener)
	go pushServer.GatewayServer.Start(listener)
}

func (pushServer *MPushServer) Stop() {
	var listener api.Listener
	pushServer.ConnectionServer.Stop(listener)
	pushServer.GatewayServer.Stop(listener)
}
