package main

import (
	client"github.com/zyl0501/go-push/sdk-server/push"
)

func main() {
	connectClient := client.ConnectClient{}
	pushClient := client.PushClient{ConnClient:connectClient}
	connectClient.Connect("localhost", 9933)
	pushClient.Start()


	defer connectClient.Close()
}
