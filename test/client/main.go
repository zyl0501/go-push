package main

import (
	"net"
	"fmt"
	"os"
	"github.com/zyl0501/go-push/test/client/config"
	"github.com/zyl0501/go-push/api/protocol"
	"encoding/json"
)

func main() {
	config.InitConfig()
	server := "localhost:9933"
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
	packet := protocol.Packet{Cmd: protocol.HANDSHAKE}
	words, _ := json.Marshal(packet)
	conn.Write([]byte(words))
	fmt.Println("send over", string(words[:]))
}
