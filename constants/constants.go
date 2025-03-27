package constants

import "time"

type CtxStrType string

const (
	UserId           CtxStrType = "user_id"
	MessageType      CtxStrType = "message_type"
	RequestTimestamp CtxStrType = "request_timestamp"
	Latency          CtxStrType = "latency"
	Path             CtxStrType = "path"
)

var (
	Consumer      string
	LocalLocation *time.Location
)