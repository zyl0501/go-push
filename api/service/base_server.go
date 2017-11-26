package service

import (
	"sync/atomic"
)

type BaseServer struct {
	BootFunc
	Started int32
}

func NewBaseServer(bootFunc BootFunc) *BaseServer {
	return &BaseServer{BootFunc: bootFunc, Started: 0}
}

func (server *BaseServer) Start(ch chan Result) {
	atomic.CompareAndSwapInt32(&server.Started, 0, 1)
	server.BootFunc.StartFunc(ch)
}

func (server *BaseServer) Stop(ch chan Result) {
	atomic.CompareAndSwapInt32(&server.Started, 1, 0)
	server.BootFunc.StopFunc(ch)
}

func (server *BaseServer) Init() {
}

func (server *BaseServer) IsRunning() (success bool) {
	return atomic.LoadInt32(&server.Started) == 1
}

type BootFunc interface {
	StartFunc(chan Result)
	StopFunc(chan Result)
}
