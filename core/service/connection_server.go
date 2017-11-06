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
	"encoding/json"
)

type ConnectionServer struct {
	server            service.BaseServer
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
}

func NewConnectionServer() (server ConnectionServer) {
	return ConnectionServer{
		server:            service.BaseServer{},
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
	}
}

func (server *ConnectionServer) Start(listener service.Listener) {
	server.server.Start(listener)
	server.listen()
}

func (server *ConnectionServer) Stop(listener service.Listener) {
	server.server.Stop(listener)
}

func (server *ConnectionServer) SyncStart() (success bool) {
	return false
}

func (server *ConnectionServer) SyncStop() (success bool) {
	return false
}

func (server *ConnectionServer) Init() {
	server.server.Init()
	server.connManager.Init()
	server.messageDispatcher.Register(protocol.HANDSHAKE, handler.HandshakeHandler{})
}

func (server *ConnectionServer) IsRunning() (success bool) {
	return false
}

func (server *ConnectionServer) listen() {
	netListen, err := net.Listen("tcp", "localhost:9933")
	if err != nil {
		log.Error(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer log.Info("server exit")
	defer netListen.Close()

	log.Info("Wait for Client")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		log.Info(conn.RemoteAddr().String(), "tcp connect success")

		server.handlerMessage(conn)
	}
}
func (server *ConnectionServer) handlerMessage(conn net.Conn) {
	serverConn := connection.NewServerConnection()
	serverConn.Init(conn)
	server.connManager.Add(serverConn)
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Error(conn.RemoteAddr().String(), "connect error:", err)
		}
	}
	packet := decodeJson(buffer[:n])
	server.messageDispatcher.OnReceive(packet, serverConn)
}

func decodeJson(content []byte) (packet protocol.Packet) {
	packet = protocol.Packet{}
	json.Unmarshal(content, &packet)
	return packet
}
