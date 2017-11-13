package main

import (
	client"github.com/zyl0501/go-push/client/push"
	"github.com/zyl0501/go-push/api/push"
)

func main() {
	connectClient := client.ConnectClient{}
	connectClient.Connect("localhost", 9933)

	pushClient :=client.PushClient{connectClient}
	pushClient.Send(push.PushContext{})

	defer connectClient.Close()
}

