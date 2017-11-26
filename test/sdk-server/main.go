package main

import (
	"time"
	"github.com/zyl0501/go-push/sdk-server/push"
	push2 "github.com/zyl0501/go-push/api/push"
	"strconv"
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/tools/config"
)

func main() {
	pushClient := push.PushClient{}
	log.Debug("%+v",config.CC.Security)
	log.Debug(config.CC.Security.PrivateKey)
	pushClient.Init()
	pushClient.Start()

	FakeBizProcess(&pushClient)

	defer pushClient.Destroy()
}

func FakeBizProcess(pushClient *push.PushClient) {
	index := 1
	for {
		pushMsg := push2.PushMsg{Content: "content_" + strconv.Itoa(index), MsgType: push2.MESSAGE, MsgId: "msgId_" + strconv.Itoa(index)}
		context := push2.PushContext{}
		context.Msg = pushMsg
		context.UserId = "user-0"
		context.Broadcast = false
		context.Timeout = 5000
		context.ACK = push2.NO_ACK

		pushClient.Send(context)
		fmt.Println("push content_" + strconv.Itoa(index))
		index++

		time.Sleep(time.Second * 5)
	}
}
