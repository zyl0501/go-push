package push

import (
	"net"
	"strconv"
	log "github.com/alecthomas/log4go"
	push "github.com/zyl0501/go-push/core/connection"
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
	conn.SetKeepAlive(true)
	client.conn = push.NewPushConnection()
	client.conn.Init(conn)
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
