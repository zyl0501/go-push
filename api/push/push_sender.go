package push

type PushSender interface {
	Send(context PushContext) (PushResult)
}
