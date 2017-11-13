package push

type PushMsg struct {
	Content string
	MsgId   string //返回使用
	MsgType byte
}

const (
	NOTIFICATION             byte = 1 + iota
	MESSAGE
	NOTIFICATION_AND_MESSAGE
)
