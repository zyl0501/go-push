package connection

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/zyl0501/go-push/api"
	log "github.com/alecthomas/log4go"
)

type ServerConnectionManager struct {
	connections cmap.ConcurrentMap
}

func NewConnectionManager() ServerConnectionManager {
	connManager := ServerConnectionManager{connections:cmap.New()}
	return connManager
}

func (manager *ServerConnectionManager) Init() {
}

func (manager *ServerConnectionManager) GetConnNum() int {
	return manager.connections.Count()
}

func (manager *ServerConnectionManager) Add(connection api.Conn) {
	manager.connections.SetIfAbsent(connection.GetId(), connection)
	log.Info("%d connections", manager.connections.Count())
}

func (manager *ServerConnectionManager) Get(id string) api.Conn {
	val, _ := manager.connections.Get(id)
	if val != nil {
		return val.(api.Conn)
	} else {
		return nil
	}
}

func (manager *ServerConnectionManager) RemoveAndClose(id string) api.Conn {
	var (
		mapVal   interface{}
		wasFound bool
	)
	cb := func(key string, val interface{}, exists bool) bool {
		mapVal = val
		wasFound = exists
		if _, ok := val.(api.Conn); ok {
			return key == id
		}
		return false
	}
	removeOk := manager.connections.RemoveCb(id, cb)
	if removeOk {
		conn := mapVal.(api.Conn)
		if conn != nil{
			conn.Close()
		}
		return conn
	} else {
		return nil
	}
}

func (manager *ServerConnectionManager) Destroy() {
	items := manager.connections.Items()
	manager.connections.IterCb(func(key string, v interface{}) {
		conn, ok := v.(api.Conn)
		if conn != nil && ok {
			conn.Close()
		}
		delete(items, key)
	})
}
