package protocol

type Packet struct {
	Cmd       byte  //命令
	cc        int16 //校验码 暂时没有用到
	flags     byte  //特性，如是否加密，是否压缩等
	sessionId int   // 会话id。客户端生成。
	lrc       byte  // 校验，纵向冗余校验。只校验head
	Body      string
}

const (
	HEARTBEAT            byte = 1 + iota //1
	HANDSHAKE                            //2
	LOGIN                                //3
	LOGOUT                               //4
	BIND                                 //5
	UNBIND                               //6
	FAST_CONNECT                         //7
	PAUSE                                //8
	ERROR                                //9
	OK                                   //10
	HTTP_PROXY                           //11
	KICK                                 //12
	GATEWAY_KICK                         //13
	PUSH                                 //14
	GATEWAY_PUSH                         //15
	NOTIFICATION                         //16
	GATEWAY_NOTIFICATION                 //17
	CHAT                                 //18
	GATEWAY_CHAT                         //19
	GROUP                                //20
	GATEWAY_GROUP                        //21
	ACK                                  //22
	NACK                                 //23
	UNKNOWN              = -1            //-1
)
