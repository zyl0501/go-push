package api

import "net"

const (
	STATUS_NEW          byte = 0
	STATUS_CONNECTED    byte = 1
	STATUS_DISCONNECTED byte = 2
)

type Conn interface {
	Init(conn net.Conn)
	GetId() string
	IsConnected() bool
	IsReadTimeout() bool
	IsWriteTimeout() bool
	UpdateLastReadTime()
	UpdateLastWriteTime()
	Close()
	GetConn() net.Conn
}

type ConnectionManager interface {
	Init()
	Add(connection Conn)
	Get(id string) Conn
	RemoveAndClose(id string) Conn
	GetConnNum() int
	Destroy()
}
