package service

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/zyl0501/go-push/api"
)

type ServerConnectionManager struct {
	connections cmap.ConcurrentMap
}

func (manager *ServerConnectionManager) Init(){
}

func (manager *ServerConnectionManager) GetConnNum() int {
	return manager.connections.Count()
}

func (manager *ServerConnectionManager) Add(connection api.Connection) {
	manager.connections.SetIfAbsent(connection.GetId(), connection)
}

func (manager *ServerConnectionManager) Get(id string) api.Connection {
	val, _ := manager.connections.Get(id)
	if val != nil {
		return val.(api.Connection)
	} else {
		return nil
	}
}

func (manager *ServerConnectionManager) RemoveAndClose(id string) api.Connection {
	var (
		mapVal interface{}
		wasFound bool

	)
	cb := func(key string, val interface{}, exists bool) bool {
		mapVal = val
		wasFound = exists
		if _, ok := val.(api.Connection); ok {
			return key == id
		}
		return false
	}
	removeOk := manager.connections.RemoveCb(id, cb)
	if removeOk {
		return mapVal.(api.Connection)
	}else{
		return nil
	}
}

func (manager *ServerConnectionManager) Destroy(){
	items := manager.connections.Items()
	manager.connections.IterCb(func(key string, v interface{}) {
		conn, ok := v.(api.Connection)
		if conn != nil && ok {
			conn.Close()
		}
		delete(items, key)
	})
}

