package main

import (
	"github.com/zyl0501/go-push/core"
)

func main() {
	server:= core.NewPushServer()
	server.Init()
	server.Start()
}
