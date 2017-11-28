package zk

import (
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/api/srd"
	log "github.com/alecthomas/log4go"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/zyl0501/go-push/tools/config"
	"encoding/json"
	"github.com/zyl0501/go-push/tools/utils"
)

type ZKServiceRegistryAndDiscovery struct {
	srd.ServiceDiscovery
	srd.ServiceRegistry
	*service.BaseServer

	conn        *zk.Conn
	nodeInfoMap map[string]map[string]srd.ServiceNode
}

func NewZKServiceRegistryAndDiscovery() *ZKServiceRegistryAndDiscovery {
	discovery := ZKServiceRegistryAndDiscovery{}
	discovery.ServiceDiscovery = &discovery
	discovery.ServiceRegistry = &discovery
	discovery.BaseServer = service.NewBaseServer(&discovery)
	discovery.nodeInfoMap = make(map[string]map[string]srd.ServiceNode)
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

func (server *ZKServiceRegistryAndDiscovery) Subscribe(watchPath string, ch chan<- srd.NodeEvent) {
	go func() {
		for {
			exist, _, existCh, err := server.conn.ExistsW(watchPath)
			if err != nil {
				return
			}
			if !exist {
				for {
					event := <-existCh
					if event.Type == zk.EventNodeCreated {
						break
					}
					_, _, existCh, err = server.conn.ExistsW(watchPath)
					if err != nil {
						return
					}
				}
			}

			nodes, _, watch, err := server.conn.ChildrenW(watchPath)
			if err != nil {
				if err == zk.ErrNoNode {
					//重新检测节点是否存在
					continue
				}
			}
			<-watch
			nodeMap := server.nodeInfoMap[watchPath]
			if nodeMap == nil {
				nodeMap = make(map[string]srd.ServiceNode)
				server.nodeInfoMap[watchPath] = nodeMap
			}
			existMap := map[string]bool{}

			// handle new add nodes
			for _, node := range nodes {
				existMap[node] = true
				if _, ok := nodeMap[node]; !ok {
					serviceNode := &srd.ServiceNode{};
					utils.FromJson(getData(server.conn, watchPath+"/"+node), serviceNode)
					nodeMap[node] = *serviceNode
					ch <- srd.NodeEvent{Node: nodeMap[node], Path: node, Type: srd.TypeServiceAdd}
				}
			}
			// handle delete nodes
			for node, _ := range nodeMap {
				if _, ok := existMap[node]; !ok {
					ch <- srd.NodeEvent{Node: nodeMap[node], Path: node, Type: srd.TypeServiceRemoved}
					delete(nodeMap, node)
				}
			}
		}
	}()
}

func (server *ZKServiceRegistryAndDiscovery) UnSubscribe(path string, c chan<- srd.NodeEvent) {
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
