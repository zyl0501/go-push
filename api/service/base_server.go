package service

type BaseServer struct {
}

func (server *BaseServer) Start(listener Listener) {
	if listener != nil {
		listener.OnSuccess("start success")
	}
}

func (server *BaseServer) Stop(listener Listener) {
	if listener != nil {
		listener.OnSuccess("stop success")
	}
}

func (server *BaseServer) SyncStart() (success bool) {
	return false
}

func (server *BaseServer) SyncStop() (success bool) {
	return false
}

func (server *BaseServer) Init() {
}

func (server *BaseServer) IsRunning() (success bool) {
	return false
}
