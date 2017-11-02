package service

type Server interface {
	Start(listener Listener)
	Stop(listener Listener)
	//CompletableFuture<Boolean> Start();
	//CompletableFuture<Boolean> Stop();
	SyncStart()(bool)
	SyncStop()(bool)
	Init()
	IsRunning()(bool)
}

type Listener interface {
	OnSuccess(args ...interface{})
	OnFailure(err error)
}
