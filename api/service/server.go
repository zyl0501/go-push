package service

type Server interface {
	Start(listener Listener)
	Stop(listener Listener)
	SyncStart()(bool)
	SyncStop()(bool)
	Init()
	IsRunning()(bool)
}

type Listener interface {
	OnSuccess(args ...interface{})
	OnFailure(err error)
}
