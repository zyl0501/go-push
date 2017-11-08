package main

import (
	"net"
	"fmt"
	"os"
	"github.com/zyl0501/go-push/test/client/config"
	"github.com/zyl0501/go-push/api/protocol"
	"reflect"
	"unsafe"
	"time"
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
	//time.Sleep(time.Second * 3)
	sender1(*conn)
	//time.Sleep(time.Second * 3)
	//sender(*conn)
}

func sender(conn net.TCPConn) {
	packet := protocol.Packet{Cmd: protocol.HANDSHAKE}
	packet.Body = stringToSliceByte("aaa")
	data := protocol.EncodePacket(packet)
	conn.Write(data[0:len(data)-2])
	time.Sleep(time.Second * 10)
	conn.Write(data[len(data)-2:])
	fmt.Println("send over",len(data), string(data))
}

func sender1(conn net.TCPConn) {
	packet := protocol.Packet{Cmd: protocol.HANDSHAKE}
	packet.Body = stringToSliceByte("aaa")
	data := protocol.EncodePacket(packet)

	data2 := make([]byte,len(data)*2)
	copy(data2[0:len(data)], data)
	copy(data2[len(data):], data)
	conn.Write(data2)
	fmt.Println("send over",len(data), string(data))
}

func stringToSliceByte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
