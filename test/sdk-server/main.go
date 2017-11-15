package main

import (
	"time"
	"github.com/zyl0501/go-push/sdk-server/push"
	push2 "github.com/zyl0501/go-push/api/push"
	"strconv"
)

func main() {
	pushClient := push.PushClient{}
	pushClient.Start()

	FakeBizProcess(&pushClient)

	defer pushClient.Destroy()
}

func FakeBizProcess(pushClient *push.PushClient) {
	index := 1
	for {
		time.Sleep(time.Second * 5)

		pushMsg := push2.PushMsg{Content: "content_" + strconv.Itoa(index), MsgType: push2.MESSAGE, MsgId: "msgId_" + strconv.Itoa(index)}
		context := push2.PushContext{}
		context.Msg = pushMsg
		context.UserId = "user-" + strconv.Itoa(index)
		context.Broadcast = false
		context.Timeout = 5000
		context.ACK = push2.NO_ACK

		pushClient.Send(context)
	}
}
