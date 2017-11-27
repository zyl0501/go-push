package zk

import (
	"github.com/zyl0501/go-push/api/service"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/zyl0501/go-push/tools/config"
	log "github.com/alecthomas/log4go"
)

type ZKClient struct {
	*service.BaseServer
	conn    *zk.Conn
	session <-chan zk.Event
}

func (server *ZKClient) Init() {
}

func (server *ZKClient) StartFunc(chan service.Result) {
	server.BaseServer.Init()
	servers := make([]string, 1)
	servers[0] = config.CC.ZK.ServerAddress
	timeout := config.CC.ZK.SessionTimeout
	conn, session, err := zk.Connect(servers, timeout)
	if err != nil {
		log.Error("zk.Connect(\"%v\", %d) error(%v)", servers, timeout, err)
		return
	}
	server.conn = conn
	server.session = session
	return
}

func (server *ZKClient) StopFunc(chan service.Result) {
}
