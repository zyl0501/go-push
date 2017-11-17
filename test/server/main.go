package main

import (
	"github.com/zyl0501/go-push/core"
	"time"
)

func main() {
	server:= core.NewPushServer()
	server.Init()
	server.Start()

	time.Sleep(time.Hour * 24)
	defer server.Stop()
}
