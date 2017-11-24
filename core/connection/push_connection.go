package connection

import (
	"net"
	"github.com/zyl0501/go-push/api"
	"time"
	"github.com/zyl0501/go-push/common/security"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/tools/config"
	"fmt"
	"github.com/zyl0501/go-push/tools/utils"
)

var (
	lId = 0
)

type PushConnection struct {
	conn          net.Conn
	status        byte
	lastReadTime  time.Time
	lastWriteTime time.Time
	id            string
	context       *api.SessionContext
}

func (conn *PushConnection) String() string {
	return fmt.Sprintf("PushConnection{conn=%v,status=%d,lastReadTime=%v,lastWriteTime=%v,id=%s,context=%v}",
		conn.conn, conn.status, conn.lastReadTime.Format(utils.FullTimeFormat), conn.lastWriteTime.Format(utils.FullTimeFormat), conn.id, conn.context)
}

func NewPushConnection() (conn *PushConnection) {
	lId++
	conn = &PushConnection{
		conn:          nil,
		status:        api.STATUS_NEW,
		lastReadTime:  time.Unix(0, 0),
		lastWriteTime: time.Unix(0, 0),
		id:            string(lId)}
	return conn
}

func (serverConn *PushConnection) Init(conn net.Conn) {
	serverConn.conn = conn
	serverConn.status = api.STATUS_CONNECTED
	serverConn.lastReadTime = time.Now()
	serverConn.lastWriteTime = time.Now()
	serverConn.context = &api.SessionContext{}
	cipher, _ := security.NewRsaCipher()
	serverConn.context.Cipher0 = cipher
	serverConn.context.Heartbeat = config.MaxHeartbeat
}

func (serverConn *PushConnection) GetId() string {
	return serverConn.id
}

func (serverConn *PushConnection) IsConnected() bool {
	return serverConn.status == api.STATUS_CONNECTED
}

func (serverConn *PushConnection) IsReadTimeout() bool {
	return time.Since(serverConn.lastReadTime) > serverConn.context.Heartbeat
}

func (serverConn *PushConnection) IsWriteTimeout() bool {
	return time.Since(serverConn.lastWriteTime) > serverConn.context.Heartbeat
}

func (serverConn *PushConnection) UpdateLastReadTime() {
	serverConn.lastReadTime = time.Now()
}
func (serverConn *PushConnection) UpdateLastWriteTime() {
	serverConn.lastWriteTime = time.Now()
}

func (serverConn *PushConnection) Close() error {
	serverConn.status = api.STATUS_DISCONNECTED
	if serverConn.conn != nil {
		return serverConn.conn.Close()
	}
	return nil
}

func (serverConn *PushConnection) GetConn() net.Conn {
	return serverConn.conn
}

func (serverConn *PushConnection) GetSessionContext() *api.SessionContext {
	return serverConn.context
}

func (serverConn *PushConnection) SetSessionContext(context *api.SessionContext) {
	serverConn.context = context
}

func (conn *PushConnection) Send(packet protocol.Packet) {
	writer := bufio.NewWriter(conn.conn)
	writer.Write(protocol.EncodePacket(packet))
	writer.Flush()
}
