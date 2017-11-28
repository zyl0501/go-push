package push

import (
	"github.com/zyl0501/go-push/api/push"
	"github.com/zyl0501/go-push/api/protocol"
	log "github.com/alecthomas/log4go"
	"io"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/common/message"
	"github.com/zyl0501/go-push/common"
	"github.com/zyl0501/go-push/sdk-server/push/handler"
	"github.com/zyl0501/go-push/core/service"
	"github.com/zyl0501/go-push/zk"
)

type PushClient struct {
	connClient        *ConnectClient
	messageDispatcher common.MessageDispatcher
	discovery zk.ZKServiceRegistryAndDiscovery
}

func (client *PushClient) Init() {
	client.messageDispatcher = common.NewMessageDispatcher()
	client.messageDispatcher.Register(protocol.OK, handler.PushUpOKHandler{})
}

func (client *PushClient) Start() {
	if client.connClient == nil {
		client.connClient = &ConnectClient{}
		client.connClient.Connect("localhost", 9934)
	}
	serverConn := client.connClient.conn
	go client.listen(serverConn)
}

func (client *PushClient) Destroy() {

}

func (client *PushClient) Send(context push.PushContext) (push.PushResult) {
	conn := client.connClient.conn

	msg := message.NewPushUpMessage0(conn)
	msg.UserId = context.UserId
	msg.Content = []byte(context.Msg.Content)
	msg.Timeout = context.Timeout
	msg.Send()
	return push.PushResult{}
}

func (client *PushClient) listen(serverConn api.Conn) {
	conn := serverConn.GetConn()

	for {
		packet, err := service.ReadPacket(conn)
		if err != nil {
			if err == io.EOF {
				log.Error("%s connect error: %v", conn.RemoteAddr().String(), err)
				break
			} else {
				log.Error("%s read error: %v", conn.RemoteAddr().String(), err)
				break
			}
		}
		client.messageDispatcher.OnReceive(*packet, serverConn)
	}
}
