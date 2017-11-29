package router

import "github.com/zyl0501/go-push/api"

const (
	LOCAL  byte = 1
	REMOTE byte = 2
)

type RouterType byte

type LocalRouter struct {
	Conn api.Conn
	Type RouterType
}

func (router *LocalRouter) GetClientType() byte{
	return router.Conn.GetSessionContext().ClientType
}

type RemoteRouter struct {
	Location ClientLocation
	Type RouterType
}
