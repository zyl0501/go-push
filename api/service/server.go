package service

type Server interface {
	Start(chan Result)
	Stop(chan Result)
	Init()
	IsRunning() (bool)
}

type Result struct {
	Success bool
	Err     error
	Args    []interface{}
}
