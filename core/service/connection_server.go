package service

import (
	"net"
	log "github.com/alecthomas/log4go"
	"os"
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/core/connection"
	"github.com/zyl0501/go-push/common"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/core/handler"
	"io"
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/core/push"
	"github.com/zyl0501/go-push/api/router"
)

type ConnectionServer struct {
	service.BaseServer
	SessionManager    *session.ReusableSessionManager
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
	pushCenter        *push.PushCenter
	routerManager     *router.LocalRouterManager
}

func NewConnectionServer(SessionManager *session.ReusableSessionManager, pushCenter *push.PushCenter, routerManager *router.LocalRouterManager) (server ConnectionServer) {
	return ConnectionServer{
		SessionManager:    SessionManager,
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
		pushCenter:        pushCenter,
		routerManager:     routerManager,
	}
}

func (server *ConnectionServer) Start(listener service.Listener) {
	server.BaseServer.Start(listener)
	server.listen()
}

func (server *ConnectionServer) Stop(listener service.Listener) {
	server.BaseServer.Stop(listener)
	server.connManager.Destroy()
}

func (server *ConnectionServer) SyncStart() (success bool) {
	return false
}

func (server *ConnectionServer) SyncStop() (success bool) {
	return false
}

func (server *ConnectionServer) Init() {
	server.BaseServer.Init()
	server.connManager.Init()
	server.messageDispatcher.Register(protocol.HANDSHAKE, handler.NewHandshakeHandler(server.SessionManager, server.connManager))
	server.messageDispatcher.Register(protocol.BIND, handler.NewBindUserHandler(server.routerManager))
}

func (server *ConnectionServer) listen() {
	netListen, err := net.Listen("tcp", "localhost:9933")
	if err != nil {
		log.Error(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer log.Info("baseServer exit")
	defer netListen.Close()

	log.Info("Wait for Client")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		log.Info("%s tcp connect success", conn.RemoteAddr().String())

		go server.handlerMessage(conn)
	}
}
func (server *ConnectionServer) handlerMessage(conn net.Conn) {
	serverConn := connection.NewPushConnection()
	serverConn.Init(conn)
	server.connManager.Add(serverConn)

	for {
		packet, err := ReadPacket(conn)
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
				server.connManager.RemoveAndClose(serverConn.GetId())
				break
			} else {
				log.Error("%s read error: %v", conn.RemoteAddr().String(), err)
				break
			}
		}
		server.messageDispatcher.OnReceive(*packet, serverConn)
	}
}
