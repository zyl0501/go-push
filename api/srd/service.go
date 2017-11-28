package srd

import (
	"strconv"
	"github.com/satori/go.uuid"
)

type ServiceRegistry interface {
	Register(node ServiceNode)
	Deregister(node ServiceNode)
}

type ServiceDiscovery interface {
	Lookup(path string) []ServiceNode
	Subscribe(path string, c chan<- ListenNode)
	UnSubscribe(path string, c chan<- ListenNode)
}

var (
	TypeServiceAdd     = 1
	TypeServiceRemoved = 2
	TypeServiceUpdated = 3
)

type ListenNode struct {
	ServiceNode
	Path string
	Type int
}

type ServiceNode struct {
	ServiceName  string
	NodeId       string
	Host         string
	Port         int
	IsPersistent bool
	Attrs map[string]interface{}
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

func (node *ServiceNode) GetAttr(key string) interface{} {
	if node.Attrs != nil {
		return node.Attrs[key]
	} else {
		return nil
	}
}
