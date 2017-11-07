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
	"encoding/json"
	"io/ioutil"
	"io"
)

type ConnectionServer struct {
	baseServer        service.BaseServer
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
}

func NewConnectionServer() (server ConnectionServer) {
	return ConnectionServer{
		baseServer:        service.BaseServer{},
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
	}
}

func (server *ConnectionServer) Start(listener service.Listener) {
	server.baseServer.Start(listener)
	server.listen()
}

func (server *ConnectionServer) Stop(listener service.Listener) {
	server.baseServer.Stop(listener)
	server.connManager.Destroy()
}

func (server *ConnectionServer) SyncStart() (success bool) {
	return false
}

func (server *ConnectionServer) SyncStop() (success bool) {
	return false
}

func (server *ConnectionServer) Init() {
	server.baseServer.Init()
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
	serverConn := connection.NewServerConnection()
	serverConn.Init(conn)
	server.connManager.Add(serverConn)

	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		if err != io.EOF {
			log.Error(conn.RemoteAddr().String(), "connect error:", err)
			server.connManager.RemoveAndClose(serverConn.GetId())
		}
	} else {
		serverConn.UpdateLastReadTime()
		packet := decodeJson(buf)
		server.messageDispatcher.OnReceive(packet, serverConn)
	}
}

func decodeJson(content []byte) (protocol.Packet) {
	packet := protocol.Packet{Cmd: protocol.OK}
	err := json.Unmarshal(content, &packet)
	if err != nil {
		log.Error("content %s parse json error %v", string(content), err)
	}
	return packet
}
