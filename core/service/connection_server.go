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
	"github.com/zyl0501/go-push/api/router"
	"github.com/zyl0501/go-push/api"
	"time"
	"context"
	"github.com/zyl0501/go-push/tools/config"
	"strconv"
)

type ConnectionServer struct {
	*service.BaseServer
	SessionManager    *session.ReusableSessionManager
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
	routerManager     *router.LocalRouterManager

	connCtx context.Context
	cancel  func()
}

func NewConnectionServer(SessionManager *session.ReusableSessionManager, routerManager *router.LocalRouterManager) (server ConnectionServer) {
	connCtx, cancel := context.WithCancel(context.Background())
	server = ConnectionServer{
		SessionManager:    SessionManager,
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
		routerManager:     routerManager,
		connCtx:           connCtx,
		cancel:            cancel,
	}
	server.BaseServer = service.NewBaseServer(&server)
	return server
}

func (server *ConnectionServer) StartFunc(ch chan service.Result) {
	if ch != nil {
		ch <- service.Result{Success: true}
	}
	server.listen()
}

func (server *ConnectionServer) StopFunc(ch chan service.Result) {
	server.connManager.Destroy()
	if ch != nil {
		ch <- service.Result{Success: true}
	}
}

func (server *ConnectionServer) Init() {
	server.BaseServer.Init()
	server.connManager.Init()
	server.messageDispatcher.Register(protocol.HANDSHAKE, handler.NewHandshakeHandler(server.SessionManager, server.connManager))
	server.messageDispatcher.Register(protocol.BIND, handler.NewBindUserHandler(server.routerManager))
	server.messageDispatcher.Register(protocol.HEARTBEAT, &handler.HeartBeatHandler{})
	server.messageDispatcher.Register(protocol.FAST_CONNECT, handler.NewFastConnectHandler(server.SessionManager))
}

func (server *ConnectionServer) listen() {
	netListen, err := net.Listen("tcp",
		config.CC.Net.ConnectServerBindIp+strconv.Itoa(config.CC.Net.ConnectServerBindPort))
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

		connCtx, cancel := context.WithCancel(context.Background())
		go server.handlerMessage(connCtx, cancel, serverConn)
		//开始检测心跳
		go server.heartbeatCheck(connCtx, cancel, serverConn)
	}
}

func (server *ConnectionServer) handlerMessage(ctx context.Context, cancel context.CancelFunc, serverConn api.Conn) {
	conn := serverConn.GetConn()
	for {
		packet, err := ReadPacket(conn)
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
			} else {
				log.Error("%s read error: %v", conn.RemoteAddr().String(), err)
			}
			break
		}
		serverConn.UpdateLastReadTime()
		server.messageDispatcher.OnReceive(*packet, serverConn)
	}
	cancel()
}

func (server *ConnectionServer) heartbeatCheck(ctx context.Context, cancel context.CancelFunc, conn api.Conn) {
	log.Info("Heartbeat: %v", conn.GetSessionContext().Heartbeat)
	timeoutTimes := 0
	for {
		select {
		case <-time.After(conn.GetSessionContext().Heartbeat):
			if conn == nil || !conn.IsConnected() {
				log.Info("heartbeat timeout times=%d, connection disconnected, conn=%v", timeoutTimes, conn);
				return;
			}
			if conn.IsReadTimeout() {
				timeoutTimes += 1
				if timeoutTimes > config.CC.Core.MaxHeartbeatTimeoutTimes {
					cancel()
					log.Info("client heartbeat timeout times=%d, do close conn=%v", timeoutTimes, conn);
					continue;
				} else {
					log.Info("client heartbeat timeout times=%d, connection=%v", timeoutTimes, conn);
				}
			} else {
				timeoutTimes = 0;
				log.Info("client heartbeat health")
			}
			continue
		case <-ctx.Done():
			ctx := conn.GetSessionContext()
			server.connManager.RemoveAndClose(conn.GetId())
			if ctx.UserId != "" {
				routerManager := server.routerManager
				routerManager.UnRegister(ctx.UserId, ctx.ClientType)
			}
			log.Info("heartbeat check cancel because of context done.")
			return
		}
	}
}
