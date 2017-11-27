package srd

import (
	"strconv"
	"github.com/zyl0501/go-push/api/service"
	"github.com/satori/go.uuid"
)

type ServiceRegistry interface {
	service.Server

	register(node ServiceNode)
	deregister(node ServiceNode)
}

type ServiceDiscovery interface {
	service.Server

	Lookup(path string) []ServiceNode
	Subscribe(path string, c chan ListenNode)
	UnSubscribe(path string, c chan ListenNode)
}

type ListenNode struct {
	ServiceNode
	Path string
}

type ServiceNode struct {
	ServiceNodeFunc
	ServiceName string
	NodeId      string
	Host        string
	Port        int
}

func (node *ServiceNode) HostAndPort() string {
	return node.Host + ":" + strconv.Itoa(node.Port)
}

func (node *ServiceNode) NodePath() string {
	if node.NodeId == "" {
		node.NodeId = uuid.NewV4().String()
	}
	return node.ServiceName + "/" + node.NodeId
}

type ServiceNodeFunc interface {
	GetAttr(string) interface{}
	IsPersistent() bool
}
