package service

type BaseServer struct {
	Server Server
}

func (server *BaseServer) Start(listener Listener) {
	listener.OnSuccess("start success")
}

func (server *BaseServer) Stop(listener Listener) {
	listener.OnSuccess("stop success")
}

func (server *BaseServer) SyncStart(success bool) {
}

func (server *BaseServer) SyncStop(success bool) {
}

func (server *BaseServer) Init() {
}

func (server *BaseServer) IsRunning(success bool) {
}
