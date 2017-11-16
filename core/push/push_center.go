package push

import (
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
)

type PushCenter struct {
	msgQueue chan message.PushUpMessage
}

func NewPushCenter() *PushCenter {
	pushCenter := PushCenter{msgQueue: make(chan message.PushUpMessage)}
	return &pushCenter
}

func (center *PushCenter) Start() {
	go func() {
		select {
		case msg, ok := <-center.msgQueue:
			if ok {
				log.Debug("receive server sdk push, now push to client")
				pushMsg := message.NewPushDownMessage0(msg.GetConnection())
				pushMsg.Content = msg.Content
				pushMsg.Send()
			}
		}
	}()
}

func (center *PushCenter) Push(message message.PushUpMessage) {
	center.msgQueue <- message
}
