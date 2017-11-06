package connection

import (
	"net"
	"github.com/zyl0501/go-push/api"
	"time"
)

var(lId = 0)

type ServerConnection struct {
	conn          net.Conn
	status        byte
	lastReadTime  int64
	lastWriteTime int64
	id            string
}

func NewServerConnection() (conn *ServerConnection) {
	lId++
	conn = &ServerConnection{conn: nil, status: api.STATUS_NEW, lastReadTime: 0, lastWriteTime: 0, id:string(lId)}
	return conn
}

func (serverConn *ServerConnection) Init(conn net.Conn) {
	serverConn.conn = conn
	serverConn.status = api.STATUS_CONNECTED
	serverConn.lastReadTime = time.Now().Unix()
}

func (serverConn *ServerConnection) GetId() string {
	return serverConn.id
}

func (serverConn *ServerConnection) IsConnected() bool {
	return serverConn.status == api.STATUS_CONNECTED
}

func (serverConn *ServerConnection) IsReadTimeout() bool {
	return serverConn.lastReadTime-time.Now().Unix() > 60*1000
}

func (serverConn *ServerConnection) IsWriteTimeout() bool {
	return serverConn.lastWriteTime-time.Now().Unix() > 60*1000
}

func (serverConn *ServerConnection) Close() {
	serverConn.status = api.STATUS_DISCONNECTED
	if serverConn.conn != nil {
		serverConn.conn.Close()
	}
}
