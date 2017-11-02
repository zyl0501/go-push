package api

import "net"

const (
	STATUS_NEW          byte = 0
	STATUS_CONNECTED    byte = 1
	STATUS_DISCONNECTED byte = 2
)

type Connection interface {
	Init(conn *net.Conn)
	GetId() string
	IsConnected() bool
	IsReadTimeout() bool
	IsWriteTimeout() bool
	Close()
}

type ConnectionManager interface {
	Init()
	Add(connection Connection)
	Get(id string) Connection
	RemoveAndClose(id string) Connection
	GetConnNum() int
	Destroy()
}