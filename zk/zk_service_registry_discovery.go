package zk

import (
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/api/srd"
	log "github.com/alecthomas/log4go"
	"github.com/samuel/go-zookeeper/zk"
)

var client *ZKClient

type ZKServiceRegistryAndDiscovery struct {
	srd.ServiceDiscovery
	srd.ServiceRegistry
	service.BaseServer
}

func (server *ZKServiceRegistryAndDiscovery) Lookup(path string) []srd.ServiceNode {

}
func (server *ZKServiceRegistryAndDiscovery) Subscribe(path string, c chan srd.ListenNode) {
	go func() {
		for {
			event := <-client.session
			log.Debug("zookeeper get a event: %s", event.State.String())
			if path == event.Path {
				data,_,_,err:=client.conn.GetW(path)
				if err != nil {
					switch event.Type {
					case zk.EventNodeCreated:
						c <-
					case zk.EventNodeDeleted:
					case zk.EventNodeDataChanged:
					case zk.EventNodeChildrenChanged:
					}
				}
			}
		}
	}()
}
func (server *ZKServiceRegistryAndDiscovery) UnSubscribe(path string, c chan srd.ListenNode) {
}

func (server *ZKServiceRegistryAndDiscovery) register(node srd.ServiceNode){
	client.conn.Create(node.NodePath())
}
func (server *ZKServiceRegistryAndDiscovery) deregister(node srd.ServiceNode){

}

func (server *ZKServiceRegistryAndDiscovery) StartFunc(ch chan service.Result) {
	client.StartFunc(ch)
}

func (server *ZKServiceRegistryAndDiscovery) StopFunc(ch chan service.Result) {
	client.StopFunc(ch)
}
