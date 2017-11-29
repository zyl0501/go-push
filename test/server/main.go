package main

import (
	"github.com/zyl0501/go-push/core"
	"time"
	log "github.com/alecthomas/log4go"
)

func main() {
	log.LoadConfiguration("")
	server:= core.NewPushServer()
	server.Init()
	server.Start()

	time.Sleep(time.Hour * 24)
	defer server.Stop()
}
