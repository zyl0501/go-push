package push

import (
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/router"
	"github.com/zyl0501/go-push/api"
)

type PushCenter struct {
	msgQueue      chan message.PushUpMessage
	routerManager *router.LocalRouterManager
}

func NewPushCenter(routerManager *router.LocalRouterManager) *PushCenter {
	pushCenter := PushCenter{msgQueue: make(chan message.PushUpMessage), routerManager: routerManager}
	return &pushCenter
}

func (center *PushCenter) Start() {
	for {
		msg, ok := <-center.msgQueue
		log.Debug(ok)
		if ok {
			log.Debug("receive server sdk push, now push to client")
			userId := msg.UserId
			//broadcast
			if userId == "" {
				routers := center.routerManager.Routers()
				for _, userRouter := range routers {
					for _, localRouter := range userRouter {
						sendPush(localRouter.Conn, msg)
					}
				}
			} else {
				clientType := msg.ClientType
				if clientType != api.CLIENT_TYPE_UNKNOWN {
					localRouter := center.routerManager.Lookup(userId, clientType)
					if localRouter != nil {
						sendPush(localRouter.Conn, msg)
					}
				} else {
					//不识别的客户端类型，推送给这个userId绑定的所有客户端
					//如果允许pc和手机同时登陆
					routers := center.routerManager.LookupAll(userId)
					if len(routers) > 0 {
						for _, localRouter := range routers {
							sendPush(localRouter.Conn, msg)
						}
					}
				}
			}
		}
	}
}

func sendPush(conn api.Conn, msg message.PushUpMessage) {
	pushMsg := message.NewPushDownMessage0(conn)
	pushMsg.Content = msg.Content
	pushMsg.Send()
}

func (center *PushCenter) Push(message message.PushUpMessage) {
	center.msgQueue <- message
}
