package push

import (
	"github.com/zyl0501/go-push/api/push"
	"github.com/zyl0501/go-push/api/protocol"
	"net"
	"time"
	"fmt"
	"reflect"
	"unsafe"
)

type PushClient struct {
	ConnClient ConnectClient
}

func (client *PushClient) Send(context push.PushContext) (push.PushResult){
	conn := client.ConnClient.conn

	packet := protocol.Packet{Cmd: protocol.HANDSHAKE}
	packet.Body = nil
	data := protocol.EncodePacket(packet)

	data2 := make([]byte,len(data)*2)
	copy(data2[0:len(data)], data)
	copy(data2[len(data):], data)
	conn.Write(data2[0:18])
	conn.Write(data2[18:])
	return push.PushResult{}
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
	conn.Write(data2[0:18])
	time.Sleep(time.Second * 10)
	conn.Write(data2[18:])
	fmt.Println("send over",len(data), string(data))

	time.Sleep(time.Second * 10)
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