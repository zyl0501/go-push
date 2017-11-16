package handler

import (
	"github.com/zyl0501/go-push/api"
	//log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/common/message"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/common/security"
	"github.com/zyl0501/go-push/core/connection"
	"github.com/zyl0501/go-push/test/client/config"
	"github.com/zyl0501/go-push/core/session"
)

type HandshakeHandler struct {
	*BaseMessageHandler
	ConnectionManager connection.ServerConnectionManager
	SessionManager    *session.ReusableSessionManager
}

func NewHandshakeHandler(SessionManager *session.ReusableSessionManager, connManager connection.ServerConnectionManager) *HandshakeHandler {
	baseHandler := &BaseMessageHandler{}
	handler := HandshakeHandler{BaseMessageHandler: baseHandler, SessionManager: SessionManager, ConnectionManager: connManager}
	handler.BaseMessageHandlerWrap = &handler
	return &handler
}

func (handler *HandshakeHandler) Decode(packet protocol.Packet, conn api.Conn) api.Message {
	return message.NewHandshakeMessage(packet, conn)
	//return &((*message.NewHandshakeMessage(packet, conn)).(api.Message))
}

func (handler *HandshakeHandler) HandleMessage(m api.Message) {
	log.Debug("HandshakeHandler HandleMessage")
	var msg *message.HandshakeMessage
	msg = m.(*message.HandshakeMessage)

	iv := msg.Iv;                                                     //AES密钥向量16位
	clientKey := msg.ClientKey;                                       //客户端随机数16位
	serverKey := security.CipherBoxIns.RandomAESKey();                //服务端随机数16位
	log.Debug("clientKey=%s, serverKey=%s", clientKey, serverKey);
	sessionKey := security.CipherBoxIns.MixKey(clientKey, serverKey); //会话密钥16位

	//1.校验客户端消息字段
	if len(msg.DeviceId) == 0 || len(iv) != security.CipherBoxIns.AesKeyLength || len(clientKey) != security.CipherBoxIns.AesKeyLength {
		errMsg := message.NewErrorMessage(msg)
		errMsg.Reason = "Param invalid"
		errMsg.Send()
		handler.ConnectionManager.RemoveAndClose(msg.GetConnection().GetId())
		log.Error("handshake failure, message=%v, conn=%v", msg, msg.GetConnection());
		return
	}
	//2.重复握手判断
	ctx := msg.GetConnection().GetSessionContext()
	if msg.DeviceId == ctx.DeviceId {
		errMsg := message.NewErrorMessage(msg)
		errMsg.Reason = "repeat handshake"
		errMsg.Send()
		log.Warn("handshake failure, repeat handshake, conn=%v", msg.GetConnection())
		return;
	}
	//3.更换会话密钥RSA=>AES(clientKey)
	ctx.Cipher0 = &security.AesCipher{clientKey, iv}
	//4.生成可复用session, 用于快速重连
	reusableSession := session.NewSession(*ctx)
	//5.计算心跳时间
	heartbeat := config.GetHeartbeat(msg.MinHeartbeat, msg.MaxHeartbeat);
	//6.响应握手成功消息
	okMsg := message.NewHandshakeOKMessage(msg.Pkt, msg.GetConnection())
	okMsg.ServerKey = serverKey
	okMsg.Heartbeat = heartbeat
	okMsg.SessionId = reusableSession.SessionId
	okMsg.ExpireTime = reusableSession.ExpireTime
	okMsg.Send()
	//7.更换会话密钥AES(clientKey)=>AES(sessionKey)
	ctx.Cipher0 = &security.AesCipher{Key: sessionKey, Iv: iv}
	//8.保存client信息到当前连接
	ctx.DeviceId = msg.DeviceId
	ctx.OsName = msg.OsName
	ctx.OsVersion = msg.OsVersion
	ctx.ClientVersion = msg.ClientVersion
	ctx.Heartbeat = heartbeat
	//ctx.ClientType =
	//9.保存可复用session到Redis, 用于快速重连
	handler.SessionManager.CacheSession(reusableSession)
}
