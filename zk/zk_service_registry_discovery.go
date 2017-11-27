package zk

import (
	"github.com/zyl0501/go-push/api/service"
	"github.com/zyl0501/go-push/api/srd"
	log "github.com/alecthomas/log4go"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/zyl0501/go-push/tools/config"
	"encoding/json"
	"github.com/zyl0501/go-push/test/go_basic/server"
)

type ZKServiceRegistryAndDiscovery struct {
	srd.ServiceDiscovery
	srd.ServiceRegistry
	*service.BaseServer

	conn    *zk.Conn
	session <-chan zk.Event
}

func (server *ZKServiceRegistryAndDiscovery) Lookup(path string) []srd.ServiceNode {

}
func (server *ZKServiceRegistryAndDiscovery) Subscribe(path string, c chan srd.ListenNode) {
	go func() {
		for {
			event := <-server.session
			log.Debug("zookeeper get a event: %s", event.State.String())
			if path == event.Path {
				data, _, _, err := server.conn.GetW(path)
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

func (server *ZKServiceRegistryAndDiscovery) register(node srd.ServiceNode) {
	data, err := json.Marshal(node)
	if err != nil {
		log.Warn("node register failure by json error. %v", node)
		return
	}
	if node.IsPersistent {
		server.registerPersist(node.NodePath(), data)
	} else {
		server.registerEphemeral(node.NodePath(), data)
	}
}
func (server *ZKServiceRegistryAndDiscovery) deregister(node srd.ServiceNode) {

}

func (server *ZKServiceRegistryAndDiscovery) StartFunc(ch chan service.Result) {
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
}

func (server *ZKServiceRegistryAndDiscovery) StopFunc(ch chan service.Result) {
}

/**
 * 持久化数据
 *
 * @param key
 * @param value
 */
func (server *ZKServiceRegistryAndDiscovery) registerPersist(key string, value []byte) {
	if server.isExisted(key) {
		server.update(key, value);
	} else {
		client.create().creatingParentsIfNeeded().withMode(CreateMode.PERSISTENT).forPath(key, value.getBytes());
	}
}

/**
 * 注册临时数据
 *
 * @param key
 * @param value
 */
func (server *ZKServiceRegistryAndDiscovery) registerEphemeral(key string, value []byte) {

}

/**
 * 更新数据
 *
 * @param key
 * @param value
 */
func (server *ZKServiceRegistryAndDiscovery) update(key string, value []byte) {

	client.inTransaction().check().forPath(key).and().setData().forPath(key, value.getBytes(Constants.UTF_8)).and().commit();
}

/**
 * 判断路径是否存在
 *
 * @param key
 * @return
 */
func (server *ZKServiceRegistryAndDiscovery) isExisted(key string) bool {
	exist, _, err := server.conn.Exists(key)
	if err != nil {
		log.Warn("check node %s exist error.", key)
		return false
	} else {
		return exist
	}
}
