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
	GetSessionContext() SessionContext
	SetSessionContext(context SessionContext)
}

type ConnectionManager interface {
	Init()
	Add(connection Conn)
	Get(id string) Conn
	RemoveAndClose(id string) Conn
	GetConnNum() int
	Destroy()
}

type SessionContext struct {
	DeviceId      string
	OsName        string
	OsVersion     string
	ClientVersion string
	UserId        string
	Tags          string
	Heartbeat     int32
	ClientType    byte
	Cipher0       Cipher
}

type Cipher interface {
	decrypt(data []byte) []byte
	encrypt(data []byte) []byte
}
