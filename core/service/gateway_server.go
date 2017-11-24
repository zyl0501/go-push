package service

import (
	"github.com/zyl0501/go-push/core/connection"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/common"
	"github.com/zyl0501/go-push/core/handler"
	"github.com/zyl0501/go-push/core/push"
	"net"
	"io"
	log "github.com/alecthomas/log4go"
	"os"
)

type GatewayServer struct {
	*service.BaseServer
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
	pushCenter        *push.PushCenter
}

func NewGatewayServer(pushCenter *push.PushCenter) *GatewayServer {
	server := GatewayServer{}
	server.BaseServer = service.NewBaseServer(&server)
	server.connManager = connection.NewConnectionManager()
	server.messageDispatcher = common.NewMessageDispatcher()
	server.pushCenter = pushCenter
	return &server
}

func (server *GatewayServer) Init() {
	server.BaseServer.Init()
	server.connManager.Init()
	server.messageDispatcher.Register(protocol.PUSH_UP, handler.NewPushUpHandler(server.pushCenter))
}

func (server *GatewayServer) StartFunc(ch chan service.Result) () {
	if ch !=nil {
		ch <- service.Result{Success: true}
	}
	server.listen()
}

func (server *GatewayServer) StopFunc(ch chan service.Result) () {
	server.connManager.Destroy()
	if ch !=nil {
		ch <- service.Result{Success: true}
	}
}

func (server *GatewayServer) listen() {
	netListen, err := net.Listen("tcp", "localhost:9934")
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

func (server *GatewayServer) handlerMessage(conn net.Conn) {
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
				continue
			}
		}
		server.messageDispatcher.OnReceive(*packet, serverConn)
	}
}
