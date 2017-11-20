package router

import (
	"github.com/orcaman/concurrent-map"
)

type LocalRouterManager struct {
	routers cmap.ConcurrentMap
}

func NewLocalRouterManager() *LocalRouterManager {
	manager := LocalRouterManager{routers: cmap.New()}
	return &manager
}

func (manager *LocalRouterManager) Register(userId string, router LocalRouter) {
	cb := func(exist bool, valueInMap interface{}, newValue interface{}) interface{} {
		var childMap map[byte]LocalRouter
		if exist {
			childMap = valueInMap.(map[byte]LocalRouter)
		} else {
			childMap = make(map[byte]LocalRouter)
		}
		childMap[router.GetClientType()] = router
		return childMap
	}
	manager.routers.Upsert(userId, router, cb)
}

func (manager *LocalRouterManager) UnRegister(userId string, clientType byte) {
	childMap, ok := manager.routers.Get(userId)
	if ok && childMap != nil {
		delete(childMap.(map[byte]LocalRouter), clientType)
	}
}

func (manager *LocalRouterManager) LookupAll(userId string) ([]LocalRouter) {
	childMap, ok := manager.routers.Get(userId)
	if ok && childMap != nil {
		childMap0 := childMap.(map[byte]LocalRouter)
		sets := make([]LocalRouter, 0, len(childMap0))
		for _, v := range childMap0 {
			sets = append(sets, v)
		}
		return sets
	} else {
		return nil
	}
}

func (manager *LocalRouterManager) Lookup(userId string, clientType byte) (*LocalRouter) {
	childMap, ok := manager.routers.Get(userId)
	if ok && childMap != nil {
		childMap0 := childMap.(map[byte]LocalRouter)
		router := childMap0[clientType]
		return &router
	} else {
		return nil
	}
}

func (manager *LocalRouterManager) Routers() (map[string](map[byte]LocalRouter)) {
	routerLen := len(manager.routers.Items())
	if routerLen <= 0 {
		return nil
	}
	var resultMap map[string](map[byte]LocalRouter) = make(map[string](map[byte]LocalRouter), routerLen)
	for key, value := range manager.routers.Items() {
		if value != nil {
			resultMap[key] = value.(map[byte]LocalRouter)
		}
	}
	return resultMap
}
