package times

import (
	"time"
	"fmt"
)

func TimeTest() {
	for {
		select {
		case t := <-time.After(1 * time.Second):
			fmt.Println(t)
			//time.Sleep(time.Second)
		}
	}
}

func TimeTest2(){
	ch := make(chan int)
	go produce(ch)
	go consumer(ch)
	time.Sleep(1 * time.Second)
}

func produce(p chan<- int) {
	for i := 0; i < 10; i++ {
		p <- i
		fmt.Println("send:", i)
	}
}
func consumer(c <-chan int) {
	for i := 0; i < 10; i++ {
		v := <-c
		fmt.Println("receive:", v)
	}
}