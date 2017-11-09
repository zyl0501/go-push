package handler

import (
	"github.com/zyl0501/go-push/api"
	//log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/common/message"
	"strings"
	log "github.com/alecthomas/log4go"
	"github.com/zyl0501/go-push/common/security"
)

type HandshakeHandler struct {
}

func (handler HandshakeHandler) Handle(packet protocol.Packet, conn api.Conn) {
	msg := message.HandshakeMessage(handler.decode(packet, conn))

	iv := msg.Iv;                  //AES密钥向量16位
	clientKey := msg.ClientKey;    //客户端随机数16位
	serverKey := make([]byte, 0);  //服务端随机数16位
	sessionKey := make([]byte, 0); //会话密钥16位

	//1.校验客户端消息字段
	if len(msg.DeviceId) == 0 || len(iv) != security.CipherBoxIns.AesKeyLength || len(clientKey) != security.CipherBoxIns.AesKeyLength {
		message.ErrorMessage{}.from(msg).setReason("Param invalid").close();
		log.Error("handshake failure, message={}, conn={}", msg, msg.GetConnection());
		return
	}
	//2.重复握手判断
	//3.更换会话密钥RSA=>AES(clientKey)
	//4.生成可复用session, 用于快速重连
	//5.计算心跳时间
	//6.响应握手成功消息
	//7.更换会话密钥AES(clientKey)=>AES(sessionKey)
	//8.保存client信息到当前连接
	//9.保存可复用session到Redis, 用于快速重连
}

func (handler HandshakeHandler) decode(packet protocol.Packet, conn api.Conn) api.Message {
	return message.HandshakeMessage{Pkt: packet, Connection: conn}
}
