package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	//children, stat, ch, err := c.ChildrenW("/")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v %+v\n", children, stat)
	//for {
	//	e := <-ch
	//	fmt.Printf("%+v\n", e)
	//}

	_, _, existCh, err := c.ExistsW("/tttt")
	fmt.Println("err:%+v", err)
	for {
		event := <-existCh
		time.Sleep(3 * time.Second)
		fmt.Println("event:%+v", event)
	}
}