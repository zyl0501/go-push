package push

import (
	"net"
	"strconv"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/common/message"
	push "github.com/zyl0501/go-push/core/connection"
	"github.com/zyl0501/go-push/common/security"
	"github.com/zyl0501/go-push/api"
)

type ConnectClient struct {
	conn api.Conn
}

func (client *ConnectClient) Connect(host string, port int) *net.TCPConn {
	server := host + ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Error("Fatal error:%s", err)
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Error("Fatal error:%s", err)
		return nil
	}
	client.conn = push.NewPushConnection()
	client.conn.Init(conn)
	client.handshake()
	return conn
}

func (client *ConnectClient) Close() {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			log.Error("Fatal error:%s", err)
		}
	}
}

func (client *ConnectClient) handshake() {
	context := client.conn.GetSessionContext()
	context.Cipher0, _ = security.NewRsaCipher()
	handshakeMsg := message.NewHandshakeMessage0(client.conn)
	handshakeMsg.DeviceId = "1111"
	handshakeMsg.OsName = "Windows"
	handshakeMsg.OsVersion = "10"
	handshakeMsg.ClientVersion = "1.0"
	handshakeMsg.Iv = security.CipherBoxIns.RandomAESIV()
	handshakeMsg.ClientKey = security.CipherBoxIns.RandomAESKey()
	handshakeMsg.MinHeartbeat = 10000
	handshakeMsg.MaxHeartbeat = 10000
	handshakeMsg.Timestamp = 0
	handshakeMsg.Send()
}
