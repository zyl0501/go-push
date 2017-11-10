package connection

import (
	"net"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/tools/config"
	"time"
)

var (
	lId = 0
)

type ServerConnection struct {
	conn          net.Conn
	status        byte
	lastReadTime  time.Time
	lastWriteTime time.Time
	id            string
	context       api.SessionContext
}

func NewServerConnection() (conn *ServerConnection) {
	lId++
	conn = &ServerConnection{
		conn:          nil,
		status:        api.STATUS_NEW,
		lastReadTime:  time.Unix(0, 0),
		lastWriteTime: time.Unix(0, 0),
		id:            string(lId)}
	return conn
}

func (serverConn *ServerConnection) Init(conn net.Conn) {
	serverConn.conn = conn
	serverConn.status = api.STATUS_CONNECTED
	serverConn.lastReadTime = time.Now()
	serverConn.lastWriteTime = time.Now()
}

func (serverConn *ServerConnection) GetId() string {
	return serverConn.id
}

func (serverConn *ServerConnection) IsConnected() bool {
	return serverConn.status == api.STATUS_CONNECTED
}

func (serverConn *ServerConnection) IsReadTimeout() bool {
	return time.Since(serverConn.lastReadTime) > config.Heartbeat
}

func (serverConn *ServerConnection) IsWriteTimeout() bool {
	return time.Since(serverConn.lastReadTime) > config.Heartbeat
}

func (serverConn *ServerConnection) UpdateLastReadTime() {
	serverConn.lastReadTime = time.Now()
}
func (serverConn *ServerConnection) UpdateLastWriteTime() {
	serverConn.lastWriteTime = time.Now()
}

func (serverConn *ServerConnection) Close() {
	serverConn.status = api.STATUS_DISCONNECTED
	if serverConn.conn != nil {
		serverConn.conn.Close()
	}
}

func (serverConn *ServerConnection) GetConn() net.Conn {
	return serverConn.conn
}

func (serverConn *ServerConnection) GetSessionContext() api.SessionContext {
	return serverConn.context
}

func (serverConn *ServerConnection) SetSessionContext(context api.SessionContext) {
	serverConn.context = context
}
