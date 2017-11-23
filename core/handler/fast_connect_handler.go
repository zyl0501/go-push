package handler

import (
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/common/message"
	"github.com/zyl0501/go-push/core/session"
	"github.com/zyl0501/go-push/test/client/config"
)

type FastConnectHandler struct {
	*BaseMessageHandler
	reusableSessionManager *session.ReusableSessionManager
}

func NewFastConnectHandler(reusableSessionManager *session.ReusableSessionManager) *FastConnectHandler {
	baseHandler := &BaseMessageHandler{}
	handler := FastConnectHandler{BaseMessageHandler: baseHandler}
	handler.BaseMessageHandlerWrap = &handler
	handler.reusableSessionManager = reusableSessionManager
	return &handler
}

func (handler *FastConnectHandler) Decode(packet protocol.Packet, conn api.Conn) api.Message {
	msg := message.NewFastConnectMessage0(conn)
	msg.Pkt = packet
	return msg
}

func (handler *FastConnectHandler) HandleMessage(m api.Message) {
	msg := m.(*message.FastConnectMessage)
	log.Info("receive fast connect message ")

	//从缓存中心查询session
	session := handler.reusableSessionManager.QuerySession(msg.SessionId)
	if session == nil {
		//1.没查到说明session已经失效了
		log.Warn("fast connect failure, session is expired, sessionId=%s, deviceId=%s, conn=%v", msg.SessionId, msg.DeviceId, msg.GetConnection())
	} else if session.Context.DeviceId != msg.DeviceId {
		//2.非法的设备, 当前设备不是上次生成session时的设备
		errMsg := message.NewErrorMessage(msg)
		errMsg.Reason = "invalid device"
		errMsg.Send()
		log.Warn("fast connect failure, not the same device, deviceId=%s, session=%v, conn=%v", msg.DeviceId, session, msg.GetConnection())
	} else {
		//3.校验成功，重新计算心跳，完成快速重连
		heartbeat := config.GetHeartbeat(msg.ExpireHeartbeat)
		session.Context.Heartbeat = heartbeat
		msg.GetConnection().SetSessionContext(session.Context)
		fastOkMsg := message.NewFastConnectOKMessage(msg.GetPacket().SessionId, msg.GetConnection())
		fastOkMsg.Heartbeat = heartbeat
		fastOkMsg.Send()
		log.Info("ast connect success, session=%v", session.Context)
	}
}
