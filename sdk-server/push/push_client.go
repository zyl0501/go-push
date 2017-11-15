package push

import (
	"github.com/zyl0501/go-push/api/push"
	"github.com/zyl0501/go-push/api/protocol"
	log "github.com/alecthomas/log4go"
	"io"
	"github.com/zyl0501/go-push/api"
)

type PushClient struct {
	connClient *ConnectClient
}

func (client *PushClient) Init() {
}

func (client *PushClient) Start() {
	if client.connClient == nil {
		client.connClient = &ConnectClient{}
		client.connClient.Connect("localhost", 9933)
	}
	serverConn := client.connClient.conn
	go client.listen(serverConn)
}

func (client *PushClient) Destroy() {

}

func (client *PushClient) handler(packet *protocol.Packet, conn *api.Conn) {
	cmd := packet.Cmd
	switch cmd {
	case protocol.OK:
		log.Info("push success: ")
	}
}

func (client *PushClient) Send(context push.PushContext) (push.PushResult) {
	conn := client.connClient.conn

	packet := protocol.Packet{Cmd: protocol.HANDSHAKE}
	packet.Body = nil
	data := protocol.EncodePacket(packet)

	data2 := make([]byte, len(data)*2)
	copy(data2[0:len(data)], data)
	copy(data2[len(data):], data)
	conn.GetConn().Write(data2[0:18])
	conn.GetConn().Write(data2[18:])
	return push.PushResult{}
}

func (client *PushClient) listen(serverConn api.Conn){
	conn := serverConn.GetConn()
	head := make([]byte, protocol.HeadLength)
	headReadLen := 0
loop:
	for {
		n, err := conn.Read(head[headReadLen:protocol.HeadLength])
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
				break loop
			}
		} else {
			if uint32(headReadLen)+uint32(n) < uint32(protocol.HeadLength) {
				headReadLen += n
			} else {
				headReadLen = 0
				packet, bodyLength := protocol.DecodePacket(head)
				readLen := 0
				body := make([]byte, bodyLength)
			bodyLoop:
				for {
					n, err := conn.Read(body[readLen: bodyLength])
					if err != nil {
						if err == io.EOF {
							log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
							break loop
						} else {
							break bodyLoop
						}
					} else {
						if uint32(readLen)+uint32(n) < bodyLength {
							readLen += n
						} else {
							packet.Body = body
							client.handler(&packet, &serverConn)
							break
						}
					}
				}
			}
		}
	}
}