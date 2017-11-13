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
)

type ConnectionServer struct {
	service.BaseServer
	connManager       connection.ServerConnectionManager
	messageDispatcher common.MessageDispatcher
}

func NewConnectionServer() (server ConnectionServer) {
	return ConnectionServer{
		connManager:       connection.NewConnectionManager(),
		messageDispatcher: common.NewMessageDispatcher(),
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
	server.messageDispatcher.Register(protocol.HANDSHAKE, handler.NewHandshakeHandler(server.connManager))
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
	serverConn := connection.NewPushConnection()
	serverConn.Init(conn)
	server.connManager.Add(serverConn)

	//var rc chan []byte
	head := make([]byte, protocol.HeadLength)
	headReadLen := 0
loop:
	for {
		n, err := conn.Read(head[headReadLen:protocol.HeadLength])
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
				server.connManager.RemoveAndClose(serverConn.GetId())
				break loop
			}
		} else {
			if uint32(headReadLen)+uint32(n) < uint32(protocol.HeadLength) {
				log.Debug("read head part %s", string(head[headReadLen:headReadLen+n]))
				headReadLen += n
			} else {
				headReadLen = 0
				log.Debug("read head complete %s", string(head))
				packet, bodyLength := protocol.DecodePacket(head)
				readLen := 0
				body := make([]byte, bodyLength)
				log.Debug("body length %d", bodyLength)
			bodyLoop:
				for {
					n, err := conn.Read(body[readLen: bodyLength])
					if err != nil {
						if err == io.EOF {
							log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
							server.connManager.RemoveAndClose(serverConn.GetId())
							break loop
						} else {
							break bodyLoop
						}
					} else {
						if uint32(readLen)+uint32(n) < bodyLength {
							log.Debug("read body part %s", string(body[readLen:readLen+n]))
							readLen += n
						} else {
							log.Debug("read body complete %s", string(body))
							packet.Body = body
							server.messageDispatcher.OnReceive(packet, serverConn)
							break
						}
					}
				}
			}
		}
	}
}
