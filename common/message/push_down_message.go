package message

type PushDownMessage struct {
	*BaseMessage

	Content []byte
}

func (msg *PushDownMessage) decodeBaseMessage(body []byte) {
	msg.Content = body
}

func (msg *PushDownMessage) encodeBaseMessage() ([]byte) {
	return msg.Content
}