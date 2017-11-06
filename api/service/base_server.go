package service

import "sync/atomic"

type BaseServer struct {
	started atomicBool
}

func (server *BaseServer) Start(listener Listener) {
	if listener != nil {
		listener.OnSuccess("start success")
	}
	server.started.compareAndSet(false, true)
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
	return server.started.get()
}

type atomicBool int32

func (b *atomicBool) get() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}
func (b *atomicBool) compareAndSet(expect bool, update bool) bool {
	var oldVal, newVal int32
	if expect {
		oldVal = 1
	} else {
		oldVal = 0
	}
	if update {
		newVal = 1
	} else {
		newVal = 0
	}
	return atomic.CompareAndSwapInt32((*int32)(b), oldVal, newVal)
}
