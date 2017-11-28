package zk

import (
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/api/srd"
	log "github.com/alecthomas/log4go"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/zyl0501/go-push/tools/config"
	"encoding/json"
	"github.com/zyl0501/go-push/tools/utils"
	"strings"
	"fmt"
)

type ZKServiceRegistryAndDiscovery struct {
	srd.ServiceDiscovery
	srd.ServiceRegistry
	*service.BaseServer

	conn *zk.Conn
}

func NewZKServiceRegistryAndDiscovery() *ZKServiceRegistryAndDiscovery {
	discovery := ZKServiceRegistryAndDiscovery{}
	discovery.ServiceDiscovery = &discovery
	discovery.ServiceRegistry = &discovery
	discovery.BaseServer = service.NewBaseServer(&discovery)
	return &discovery
}

func (server *ZKServiceRegistryAndDiscovery) Init() {}

func (server *ZKServiceRegistryAndDiscovery) StartFunc(ch chan service.Result) {
	server.BaseServer.Init()
	servers := make([]string, 1)
	servers[0] = config.CC.ZK.ServerAddress
	timeout := config.CC.ZK.SessionTimeout
	conn, _ := connect(servers, timeout)
	server.conn = conn
}

func (server *ZKServiceRegistryAndDiscovery) StopFunc(ch chan service.Result) {
	server.conn.Close()
}

func (server *ZKServiceRegistryAndDiscovery) Lookup(fpath string) []srd.ServiceNode {
	paths := getNodes(server.conn, fpath)
	if len(paths) == 0 {
		return nil
	}
	nodes := make([]srd.ServiceNode, len(paths))
	for index, path := range paths {
		if path != "" {
			data := getData(server.conn, fpath+"/"+path)
			if data != nil {
				log.Debug("lookup.data=%s", string(data))
				node := &srd.ServiceNode{}
				utils.FromJson(data, node)
				nodes[index] = *node
			}
		}
	}
	return nodes
}

func (server *ZKServiceRegistryAndDiscovery) Subscribe(path string, c chan<- srd.ListenNode) {
	go func() {
		for {
			children, ch, err := getNodesW(server.conn, path)
			if err != nil {
				fmt.Printf("%+v", err)
				return
			}
			fmt.Printf("%+v\n", children)

			event := <-ch
			log.Debug("zookeeper get a event: %v, path=%s", event, event.Path)
			if strings.HasPrefix(event.Path, path) {
				data := getData(server.conn, event.Path)
				if len(data) <= 0 {
					log.Debug("zk get data error. path=%s", event.Path)
				} else {
					node := srd.ListenNode{}
					utils.FromJson(data, node)
					switch event.Type {
					case zk.EventNodeCreated:
						node.Type = srd.TypeServiceAdd
						c <- node
					case zk.EventNodeDeleted:
						node.Type = srd.TypeServiceRemoved
						c <- node
					case zk.EventNodeDataChanged:
						node.Type = srd.TypeServiceUpdated
						c <- node
					case zk.EventNodeChildrenChanged:
						node.Type = srd.TypeServiceUpdated
						c <- node
					}
				}
			}
		}
	}()
}

func (server *ZKServiceRegistryAndDiscovery) UnSubscribe(path string, c chan<- srd.ListenNode) {
}

func (server *ZKServiceRegistryAndDiscovery) Register(node srd.ServiceNode) {
	data, err := json.Marshal(node)
	if err != nil {
		log.Debug("node register failure by json error. %v", node)
		return
	} else {
		log.Debug("node register json= %s", string(data))
	}
	if node.IsPersistent {
		registerPersist(server.conn, node.NodePath(), data)
	} else {
		registerEphemeral(server.conn, node.NodePath(), data, true)
	}
}
func (server *ZKServiceRegistryAndDiscovery) Deregister(node srd.ServiceNode) {
	remove(server.conn, node.NodePath())
}
