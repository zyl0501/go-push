package push

type PushResult struct {
	ResultCode int
	UserId string
}

const(
	CODE_SUCCESS int = 1+iota
	CODE_FAILURE
	CODE_OFFLINE
	CODE_TIMEOUT
)