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
	"github.com/zyl0501/go-push/api"
	"time"
	"context"
	"github.com/zyl0501/go-push/tools/config"
)

type ConnectionServer struct {
	service.BaseServer
	SessionManager    *session.ReusableSessionManager
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
	pushCenter        *push.PushCenter
	routerManager     *router.LocalRouterManager

	connCtx context.Context
	cancel  func()
}

func NewConnectionServer(SessionManager *session.ReusableSessionManager, pushCenter *push.PushCenter,
	routerManager *router.LocalRouterManager) (server ConnectionServer) {
	connCtx, cancel := context.WithCancel(context.Background())
	return ConnectionServer{
		SessionManager:    SessionManager,
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
		pushCenter:        pushCenter,
		routerManager:     routerManager,
		connCtx:           connCtx,
		cancel:            cancel,
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
	server.messageDispatcher.Register(protocol.HEARTBEAT, &handler.HeartBeatHandler{})
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

		serverConn := connection.NewPushConnection()
		serverConn.Init(conn)
		server.connManager.Add(serverConn)

		go server.handlerMessage(serverConn)
	}
}

func (server *ConnectionServer) handlerMessage(serverConn api.Conn) {
	conn := serverConn.GetConn()
	connCtx, cancel := context.WithCancel(context.Background())
	for {
		//开始检测心跳
		go server.heartbeatCheck(connCtx, serverConn, 0)
		packet, err := ReadPacket(conn)
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
			} else {
				log.Error("%s read error: %v", conn.RemoteAddr().String(), err)
			}
			ctx := serverConn.GetSessionContext()
			server.connManager.RemoveAndClose(serverConn.GetId())
			if ctx.UserId != "" {
				routerManager := server.routerManager
				routerManager.UnRegister(ctx.UserId, ctx.ClientType)
			}
			break
		}
		server.messageDispatcher.OnReceive(*packet, serverConn)
	}
	cancel()
}

func (server *ConnectionServer) heartbeatCheck(ctx context.Context, conn api.Conn, timeoutTimes int) {
	select {
	case t:= <-time.After(conn.GetSessionContext().Heartbeat):
		log.Info("time: ", t)
		if conn == nil || !conn.IsConnected() {
			log.Info("heartbeat timeout times=%d, connection disconnected, conn=%v", timeoutTimes, conn);
			return;
		}
		if conn.IsReadTimeout() {
			timeoutTimes += 1
			if timeoutTimes > config.MaxHeartbeatTimeoutTimes {
				server.connManager.RemoveAndClose(conn.GetId())
				log.Info("client heartbeat timeout times=%d, do close conn=%v", timeoutTimes, conn);
				return;
			} else {
				log.Info("client heartbeat timeout times=%d, connection=%v", timeoutTimes, conn);
			}
		} else {
			timeoutTimes = 0;
		}
		go server.heartbeatCheck(ctx, conn, timeoutTimes)
		return
	case <-ctx.Done():
		log.Info("heartbeat check cancel because of context done.")
		return
	}
}
