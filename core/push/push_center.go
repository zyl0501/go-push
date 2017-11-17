package push

import (
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/router"
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
			//clientType := msg.ClientType
			//localRouter := center.routerManager.Lookup(userId, clientType)
			//if localRouter != nil {
			//	pushMsg := message.NewPushDownMessage0(localRouter.Conn)
			//	pushMsg.Content = msg.Content
			//	pushMsg.Send()
			//}

			routers := center.routerManager.LookupAll(userId)
			if len(routers) > 0 {
				for _, localRouter := range routers {
					pushMsg := message.NewPushDownMessage0(localRouter.Conn)
					pushMsg.Content = msg.Content
					pushMsg.Send()
				}
			}
		}
	}
}

func (center *PushCenter) Push(message message.PushUpMessage) {
	center.msgQueue <- message
}
