package main

import (
	"sync"
	"github.com/zyl0501/go-push/api/srd"
	"fmt"
	"github.com/zyl0501/go-push/zk"
	"github.com/zyl0501/go-push/tools/config"
	"time"
)

var PathTest = "/test"

func main() {
	var wg sync.WaitGroup
	config.CC.ZK.ServerAddress = "127.0.0.1:2181"
	config.CC.ZK.SessionTimeout = time.Second * 10
	discovery := zk.NewZKServiceRegistryAndDiscovery()
	discovery.Init()
	discovery.Start(nil)

	wg.Add(1)
	c := make(chan srd.NodeEvent)
	discovery.Subscribe(PathTest, c)

	go func() {
		c := time.Tick(10 * time.Second)
		var i =1
		for{
			<-c
			n := srd.ServiceNode{}
			n.ServiceName = PathTest
			n.Host = "127.0.0.1"
			n.Port = i
			discovery.Register(n)
			i++
		}
	}()


	go func() {
		time.Sleep(7 * time.Second)
		nodes := discovery.Lookup(PathTest)
		for _, node := range nodes {
			fmt.Println("=================node: %+v", node)
		}
	}()
	for {
		node := <-c
		fmt.Println("listen node: ", node)
	}
	wg.Wait()
}
