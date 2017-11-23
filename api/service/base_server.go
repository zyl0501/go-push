package service

import (
	"sync/atomic"
)

type BaseServer struct {
	BootFunc
	started int32
}

func (server *BaseServer) Start(ch chan Result) {
	atomic.CompareAndSwapInt32(&server.started,0,1)
	server.BootFunc.StartFunc(ch)
}

func (server *BaseServer) Stop(ch chan Result) {
	atomic.CompareAndSwapInt32(&server.started,1,0)
	server.BootFunc.StopFunc(ch)
}

func (server *BaseServer) Init() {
}

func (server *BaseServer) IsRunning() (success bool) {
	return atomic.LoadInt32(&server.started) == 1
}

type BootFunc interface{
	StartFunc(chan Result)
	StopFunc(chan Result)
}
