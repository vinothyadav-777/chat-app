package constants

type CtxStrType string

const (
	UserId           CtxStrType = "user_id"
	MessageType      CtxStrType = "message_type"
	RequestTimestamp CtxStrType = "request_timestamp"
	Latency          CtxStrType = "latency"
	Path             CtxStrType = "path"
)
