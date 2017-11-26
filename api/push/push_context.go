package push

type PushContext struct {
	//content   []byte   //待推送的内容
	Msg       PushMsg  //待推送的消息
	UserId    string   //目标用户
	UserIds   []string //目标用户，批量
	ACK       byte     //消息ack模式
	Timeout   int64    //推送超时时间
	Broadcast bool     //全网广播在线用户
	Tags      []string //用户标签过滤,目前只有include, 后续会增加exclude
}

const (
	NO_ACK   byte = iota
	AUTO_ACK
	BIZ_ACK
)
