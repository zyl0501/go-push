package router

type RemoteRouterManager struct {
}

func (manager *RemoteRouterManager) Register(userId string, router RemoteRouter) *RemoteRouter {

}

func (manager *RemoteRouterManager) UnRegister(userId string, clientType byte) bool {

}

func (manager *RemoteRouterManager) LookupAll(userId string) ([]RemoteRouter) {

}

func (manager *RemoteRouterManager) Lookup(userId string, clientType byte) (*RemoteRouter) {

}
