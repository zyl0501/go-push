package router

import "github.com/zyl0501/go-push/api"

const (
	LOCAL  byte = 1
	REMOTE byte = 2
)

type LocalRouter struct {
	Conn api.Conn
	RouterType byte
}

func (router *LocalRouter) GetClientType() byte{
	return router.Conn.GetSessionContext().ClientType
}

type RemoteRouter struct {
	Location ClientLocation
	RouterType byte
}
