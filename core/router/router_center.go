package router

import "github.com/zyl0501/go-push/api/router"

type RouterCenter struct {
	LocalManager *router.LocalRouterManager
	RemoteManager *router.RemoteRouterManager
}