package main

import (
	"net"
	"fmt"
	"os"
)


func main() {
	InitConfig()
	server := "localhost:9932"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error:%s", err)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error:%s", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer fmt.Println("connect closed")
	fmt.Println("connect success")
	sender(*conn)
}

func sender(conn net.TCPConn) {
	words := "test tcp conn"
	conn.Write([]byte(words))
	fmt.Println("send over")
}
